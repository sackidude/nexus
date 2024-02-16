package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	main := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/index.html"))
		tmpl.Execute(w, nil)
	}
	http.HandleFunc("/", main)

	// HTMX triggers
	imageRequest := func(w http.ResponseWriter, r *http.Request) {
		htmlStr := fmt.Sprintf("image nr. %d", 1)
		tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
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
