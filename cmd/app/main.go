package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ebuckley/bff/pkg/bff"
	"github.com/ebuckley/bff/pkg/server"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	app := bff.New()
	err := app.RegisterAction("upload a file", func(ctx context.Context, io *bff.Io) error {
		f, err := io.Input.File("Upload a file")
		if err != nil {
			return err
		}
		io.Display.Metadata([]bff.MetadataItem{
			{Label: "File Name", Value: f[0]},
		})
		return nil
	}, bff.WithSlug("upload_file"))
	if err != nil {
		panic(err)
	}

	err = app.RegisterAction("user_profile", func(ctx context.Context, io *bff.Io) error {
		email, err := io.Input.Email("Enter your email")
		if err != nil {
			return err
		}

		age, err := io.Input.Slider("Select your age", 18, 100)
		if err != nil {
			return err
		}

		birthdate, err := io.Input.Date("Enter your birthdate")
		if err != nil {
			return err
		}

		bio, err := io.Input.TextArea("Tell us more about you")
		if err != nil {
			return err
		}

		website, err := io.Input.URL("Enter your website")
		if err != nil {
			return err
		}

		availableTime, err := io.Input.Time("Select your available time")
		if err != nil {
			return err
		}

		avatar, err := io.Input.File("Upload your avatar")
		if err != nil {
			return err
		}

		// just display it as metadata for now
		io.Display.Metadata([]bff.MetadataItem{
			{Label: "Email", Value: email},
			{Label: "Age", Value: fmt.Sprint(age)},
			{Label: "Birthdate", Value: fmt.Sprint(birthdate)},
			{Label: "Bio", Value: bio},
			{Label: "Website", Value: website},
			{Label: "Available Time", Value: fmt.Sprint(availableTime)},
			{Label: "Avatar", Value: fmt.Sprint(avatar)},
		}, bff.WithMetadataLayout("table"))

		return nil
	})
	if err != nil {
		panic(err)
	}

	err = app.RegisterAction("hello", func(ctx context.Context, io *bff.Io) error {

		io.Display.Heading("Hello World!", 1)
		io.Display.Image("https://media.giphy.com/media/26ybw6AltpBRmyS76/giphy.gif", "gopher", "medium")

		io.Display.Link("Visit Go's website", "https://golang.org", bff.WithLinkType("primary"))

		io.Display.Html("<p>This is <strong>HTML</strong> content rendered directly.</p>")

		io.Display.Code(`
package main

import "fmt"

func main() {
    fmt.Println("Hello, BFF!")
}
    `, "go")

		// New Metadata component
		io.Display.Metadata([]bff.MetadataItem{
			{Label: "Framework", Value: "BFF"},
			{Label: "Language", Value: "Go"},
			{Label: "Purpose", Value: "Backend for Frontend"},
		}, bff.WithMetadataLayout("card"))

		// Existing input component
		name, err := io.Input.Text("What is your name?")
		if err != nil {
			return err
		}

		// Existing heading component
		io.Display.Heading("Hello, "+name, 1)

		return nil
	})
	if err != nil {
		panic(err)
	}

	err = app.RegisterAction("launch nukes", launchNukes, bff.WithSlug("nuke"))
	if err != nil {
		panic(err)
	}

	s := server.NewServer(app, server.Prefix("/backend"))
	slog.Info("starting server on :8181")

	err = http.ListenAndServe(":8181", logger(s))
	if err != nil {
		panic(err)
	}
}

// logger is a basic http request logger
func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("request", "method", r.Method, "path", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func launchNukes(ctx context.Context, io *bff.Io) error {
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
	killCivvies, err := io.Input.Boolean("spare civilians?")
	if err != nil {
		return err
	}

	countDown, err := io.Input.Number("How many seconds until launch?")
	if err != nil {
		return err
	}
	for i := countDown; i > 0; i-- {
		io.Display.Heading(fmt.Sprintf("Launching it in %ds", i), 1)
		time.Sleep(1 * time.Second)
	}

	io.Display.Heading("Great job destroying "+city, 1)
	if !killCivvies {
		io.Display.Heading("You spared the civilians, you are a good person", 2)
	}
	io.Display.Markdown(`
# the Benifits of ethical nuke launching

- An ethical hacker would never make a nuke <city> button'
- Instead you should provide a <city> parameter, henceforth unloading your nuke responsibility to someone else
- This is the way of the ethical hacker
`)

	return nil
}
