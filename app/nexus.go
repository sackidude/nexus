package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

// TODO: replace all _ with real errorhandling and fallback options

var imgNum = 1

func getNewImageData() map[string]string {
	// TODO: Request data from db

	directoryName := "trial-3"
	filename := fmt.Sprintf("%d.jpg", imgNum)
	imgNum++
	time := "16/2-24"
	//
	return map[string]string{
		"trial": directoryName,
		"image": filename,
		"time":  time,
		"path":  fmt.Sprintf("images/%s/%s", directoryName, filename),
	}
}

func main() {
	// Basic static hosting
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// HTMX triggers
	// Preload fragments
	fetcherHTML, _ := os.ReadFile("fragments/fetcher/main.html")

	// HTTP Requests
	imageRequest := func(w http.ResponseWriter, r *http.Request) {
		// Request the relevant information to db.
		// Template html fragment with the information from db
		tmpl, _ := template.New("t").Parse(string(fetcherHTML))
		imageData := getNewImageData()
		tmpl.Execute(w, imageData)
	}
	http.HandleFunc("/image-request", imageRequest)

	dataViewer := func(w http.ResponseWriter, r *http.Request) {
		htmlStr := fmt.Sprintf("Here are the graphs %d", 1)
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/data-view", dataViewer)

	// HTTP Posts
	imageDataRetrival := func(w http.ResponseWriter, r *http.Request) {
		pxHeightString := r.Header.Get("pxHeight")
		pxHeight, _ := strconv.ParseFloat(pxHeightString, 64)
		fmt.Printf("User image data received with pxHeight: %.1f", pxHeight)

		imageData := getNewImageData()
		tmpl, _ := template.New("t").Parse(string(fetcherHTML))
		tmpl.Execute(w, imageData) //TODO: Fix this temporary solution
	}
	http.HandleFunc("/user-image-data", imageDataRetrival)

	http.ListenAndServe("localhost:8080", nil)
}
