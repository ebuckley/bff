package server

import (
	"bff/pkg/bff"
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log/slog"
	"net/http"
)

type Server struct {
	BFF     *bff.BFF
	session map[string]*websocket.Conn
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// is this a websocket upgrade request?
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "expected websocket connection", http.StatusBadRequest)
		return
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true, OriginPatterns: []string{"*"}})
	if err != nil {
		http.Error(w, "could not open websocket connection", http.StatusBadRequest)
	}
	defer c.CloseNow()

	// Set the context as needed. Use of r.Context() is not recommended
	// to avoid surprising behavior (see http.Hijacker).
	ctx := context.Background()
	// TODO upgrade to read more than just the ping
	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		slog.Error("failed to read bff.Message: ", "err", err)
		c.Close(websocket.StatusInternalError, "failed to read bff.Message")
		return
	}
	// this should be some kind of session registration.
	// then we can drive the state somehow -- maybe just return with a list of actions for now?

	pages := s.BFF.GetPages()
	actions := s.BFF.GetActions()
	err = send(ctx, c, bff.Message{"pages", pages})

	if err != nil {
		slog.Error("failed to write pages: ", "err", err)
		c.Close(websocket.StatusInternalError, "failed to write pages")

		return
	}
	err = send(ctx, c, bff.Message{"actions", actions})
	if err != nil {
		slog.Error("failed to write actions: ", "err", err)
		c.Close(websocket.StatusInternalError, "failed to write actions")
		return
	}

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
