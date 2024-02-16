package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

// TODO: replace all _ with real errorhandling and fallback options

func main() {
	// Basic stating hosting
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// HTMX triggers
	// Preload fragments
	fetcherHTML, _ := os.ReadFile("fragments/fetcher/main.html")

	imageRequest := func(w http.ResponseWriter, r *http.Request) {
		// Request the relevant information to db.
		directoryName := "trial-3"
		filename := "1.jpg"
		time := "16/2-24"
		// Template html fragment with the information from db
		tmpl, _ := template.New("t").Parse(string(fetcherHTML))
		tmpl.Execute(w, map[string]string{
			"trial": directoryName,
			"image": filename,
			"time":  time,
			"path":  fmt.Sprintf("images/%s/%s", directoryName, filename),
		})
	}
	http.HandleFunc("/image-request", imageRequest)

	dataViewer := func(w http.ResponseWriter, r *http.Request) {
		htmlStr := fmt.Sprintf("Here are the graphs %d", 1)
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/data-view", dataViewer)

	http.ListenAndServe("localhost:8080", nil)
}
