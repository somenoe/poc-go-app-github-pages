package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const REPOSITORY_NAME = "poc-go-app-github-pages"

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
	resp, err := http.Get("/" + REPOSITORY_NAME + "/time.json")
	if err != nil {
		h.buildTimestamp = "Error fetching time"
		log.Println("Error fetching time:", err)
		return
	}
	defer resp.Body.Close()

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)

	var timeData struct {
		Time string `json:"time"`
	}

	err = json.Unmarshal(body, &timeData)
	if err != nil {
		h.buildTimestamp = "Error parsing time"
		log.Println("Error parsing JSON:", err)
		return
	}

	log.Println("Build time:", timeData.Time)
	h.buildTimestamp = timeData.Time
}

func main() {
	var currentTime = time.Now().Format("2006-01-02 15:04:05")

	app.Route("/", func() app.Composer {
		return &hello{}
	})

	app.RunWhenOnBrowser()

	// Create time.json file
	timeJSON := fmt.Sprintf(`{"time":"%s"}`, currentTime)
	err := os.WriteFile("docs/time.json", []byte(timeJSON), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = app.GenerateStaticWebsite("docs", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Resources:   app.GitHubPages(REPOSITORY_NAME),
	})

	if err != nil {
		log.Fatal(err)
	}
}
