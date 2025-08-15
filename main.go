package main

import (
	"log"
	"net/http"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type hello struct {
	app.Compo
	buildTimestamp string
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Hello World!"),
		app.P().Text("Current time: "+h.buildTimestamp),
	)
}

func (h *hello) OnMount(ctx app.Context) {
	resp, err := http.Get("/time")
	if err != nil {
		h.buildTimestamp = "Error fetching time"
		log.Println("Error fetching time:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	timeText := string(body)

	log.Println("Build time:", timeText)
	h.buildTimestamp = timeText
}

func main() {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	app.Route("/", func() app.Composer {
		return &hello{}
	})

	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	http.Handle("/time", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Current time: " + currentTime))
	}))

	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
