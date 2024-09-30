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
		slog.Error("failed to read message: ", "err", err)
		c.Close(websocket.StatusInternalError, "failed to read message")
		return
	}
	// this should be some kind of session registration.
	// then we can drive the state somehow -- maybe just return with a list of actions for now?

	pages := s.BFF.GetPages()
	actions := s.BFF.GetActions()
	err = send(ctx, c, "pages", pages)

	if err != nil {
		slog.Error("failed to write pages: ", "err", err)
		c.Close(websocket.StatusInternalError, "failed to write pages")

		return
	}
	err = send(ctx, c, "actions", actions)
	if err != nil {
		slog.Error("failed to write actions: ", "err", err)
		c.Close(websocket.StatusInternalError, "failed to write actions")
		return
	}
	// now wait forever for more actions from the user
	for {
		var v = map[string]any
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			slog.Error("failed to read action selection: ", "err", err)
			c.Close(websocket.StatusInternalError, "failed to read message")
			return
		}

		slog.Info("received message: ", "msg", v)
		if v["type"].(string) == "start" {
			// dispatch the action and wait for things to happen
			actionName := v["action"].(string)
			_, err = s.BFF.ExecuteAction(ctx, actionName, nil)
			if err != nil {
				slog.Error("failed to execute action: ", "err", err)
				c.Close(websocket.StatusInternalError, "failed to execute action")
				return
			}

		}
	}
	c.Close(websocket.StatusNormalClosure, "")
}

func send(ctx context.Context, c *websocket.Conn, msgType string, payload any) error {
	slog.Debug("sending message: ", "type", msgType, "payload", payload)
	type content struct {
		Type string
		Data any
	}
	err := wsjson.Write(ctx, c, content{Type: msgType, Data: payload})
	if err != nil {
		return err
	}
	return nil
}
