package main

import (
	"bff/pkg/bff"
	"bff/pkg/server"
	"context"
	"log/slog"
	"net/http"
)

func main() {
	app := bff.New("development")
	err := app.RegisterAction("hello", func(ctx context.Context, io *bff.Io, params map[string]any) (any, error) {
		io.Display.Heading("Hello World!", 1)
		return nil, nil
	})
	if err != nil {
		panic(err)
	}
	s := server.Server{BFF: app}
	slog.Info("starting server on :8181")
	err = http.ListenAndServe(":8181", &s)
	if err != nil {
		panic(err)
	}
}
