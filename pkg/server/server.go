package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/ebuckley/bff/pkg/bff"
)

// TODO make configurable somehow, when turned on this flag means we will proxy to the vite app
var development string

type Server struct {
	BFF           *bff.BFF
	mux           http.Handler
	assets        http.Handler
	reactIndex    http.Handler
	handlerPrefix string
}

// NewServer creates a new server with the given BFF instance and handler prefix, the handler prefix is an optional
// string I.E `/dashboard` that will be prepended to all routes, you do not need to include a trailing slash on the prefix
func NewServer(bff *bff.BFF, opts ...Serveropts) *Server {
	s := &Server{BFF: bff}
	for _, opt := range opts {
		opt(s)
	}
	s.setup()
	return s
}

func (s *Server) setup() http.Handler {

	mux := http.NewServeMux()
	if s.handlerPrefix != "" {
		slog.Info("setup bff server with a prefix", "prefix", s.handlerPrefix)
	}
	if development == "true" {
		slog.Info("development mode enabled, proxying to vite")
	}

	// basic URL scheme is:
	// / -> index.html
	// /a/{a} -> action (a react app)
	// /a/{a}/ws -> websocket for action to do stuff
	s.assets = s.makeStaticServer()
	s.reactIndex = serveReactIndex(s.handlerPrefix)

	mux.HandleFunc(s.handlerPrefix+"/", s.index)
	mux.Handle(s.handlerPrefix+"/a/{action}", s.reactIndex)
	mux.HandleFunc(s.handlerPrefix+"/a/{action}/ws", s.handleAction)

	s.mux = mux
	return mux
}
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != s.handlerPrefix+"/" {
		// fallback to the static file server
		s.assets.ServeHTTP(w, r)
		return
	}
	state := struct {
		Prefix  string
		Heading string
		Actions []*bff.Action
	}{
		Heading: "Actions",
		Prefix:  s.handlerPrefix,
		Actions: s.BFF.GetActions(),
	}
	err := index.Execute(w, state)
	if err != nil {
		slog.Error("failed to render index", "err", err)
		http.Error(w, "failed to render index", http.StatusInternalServerError)
	}
}
func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	// is this a websocket upgrade request?
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "expected websocket connection", http.StatusBadRequest)
		return
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true, OriginPatterns: []string{"*"}})
	if err != nil {
		http.Error(w, "could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.CloseNow()
		if err != nil {
			slog.Error("failed to close websocket connection", "err", err)
		}
	}(c)

	// Set the context as needed. Use of r.Context() is not recommended
	// to avoid surprising behavior (see http.Hijacker).
	ctx := context.Background()

	// this should be some kind of session registration.
	// then we can drive the state somehow -- maybe just return with a list of actions for now?

	// now wait forever for more actions from the user
	input := make(chan bff.Message)
	output := make(chan bff.Message, 1)

	go s.BFF.Loop(ctx, input, output)

	//starts a thread for output
	go func(ctx context.Context, output chan bff.Message, c *websocket.Conn) {
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-output:
				slog.Debug("sending anotha bff.Message: ", "type", v.Type, "payload", v.Data)
				err = send(ctx, c, v)
				if err != nil {
					slog.Error("failed to write display: ", "err", err)
					c.Close(websocket.StatusInternalError, "failed to write display")
					return
				}
			}
		}
	}(ctx, output, c)

	for {
		var v bff.Message
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			slog.Error("failed to read from looped reader: ", "err", err)
			c.Close(websocket.StatusInternalError, "failed to read bff.Message")
			return
		}
		select {
		case <-ctx.Done():
			slog.Debug("closing connection")
			c.Close(websocket.StatusNormalClosure, "")
			return
		case input <- v:
			slog.Debug("received bff.Message: ", "type", v.Type, "payload", v.Data)
			// carry on reading
		}
	}
}

func send(ctx context.Context, c *websocket.Conn, m bff.Message) error {
	slog.Debug("sending bff.Message: ", "type", m.Type, "payload", m.Data)
	err := wsjson.Write(ctx, c, m)
	if err != nil {
		return err
	}
	return nil
}

type Serveropts func(s *Server)

// Prefix lets you set a handler prefix I.E "/admin" for the server, this is useful if you want to run the server under a subpath
func Prefix(prefix string) Serveropts {
	return func(s *Server) {
		s.handlerPrefix = prefix
	}
}
