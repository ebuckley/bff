package main

import (
	"bff/pkg/bff"
	"bff/pkg/server"
	"context"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	app := bff.New("development")
	err := app.RegisterAction("hello", func(ctx context.Context, io *bff.Io) error {
		io.Display.Heading("Hello World!", 1)
		name, err := io.Input.Text("What is your name?")
		if err != nil {
			return err
		}
		io.Display.Heading("Hello, "+name, 1)
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = app.RegisterAction("launch nukes", func(ctx context.Context, io *bff.Io) error {
		io.Display.Heading("Read to launch some nukes?!", 1)
		io.Display.Markdown(`
## Don't worry this is just a simulation

In this example you will see a few cool things like Yes/No booleans, text inputs markdown outputs and the way that this all happens in realtime.
`)
		confirm, err := io.Input.Boolean("Are you sure you want to launch the nuke?")
		if err != nil {
			return err
		}
		if !confirm {
			io.Display.Heading("Nuke launch aborted, you are a good person", 1)
			return nil
		}
		io.Display.Heading("Great! Let's plan a nuke launch!", 1)
		city, err := io.Input.Text("What city would you like to destroy?")
		if err != nil {
			return err
		}
		civiliansTarget, err := io.Input.Boolean("spare civilians?")
		if err != nil {
			return err
		}
		io.Display.Heading("Launching it in 10s", 1)
		time.Sleep(5 * time.Second)
		io.Display.Heading("Launching it in 5", 1)
		time.Sleep(1 * time.Second)
		io.Display.Heading("Launching in 4", 1)
		time.Sleep(1 * time.Second)
		io.Display.Heading("Launching in 3", 1)
		time.Sleep(1 * time.Second)
		io.Display.Heading("Launching in 2", 1)
		time.Sleep(1 * time.Second)
		io.Display.Heading("Launching in 1", 1)
		time.Sleep(1 * time.Second)
		io.Display.Heading("Great job destroying "+city, 1)
		if !civiliansTarget {
			io.Display.Heading("You spared the civilians, you are a good person", 2)
		}
		io.Display.Markdown(`
# the Benifits of ethical nuke launching

- An ethical hacker would never make a nuke <city> button'
- Instead you should provide a <city> parameter, henceforth unloading your nuke responsibility to someone else
- This is the way of the ethical hacker
`)

		return nil
	})
	s := server.Server{BFF: app}
	slog.Info("starting server on :8181")
	err = http.ListenAndServe(":8181", &s)
	if err != nil {
		panic(err)
	}
}
