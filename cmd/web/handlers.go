package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// HTTP GET
func ImageRequest(w http.ResponseWriter, r *http.Request) {
	tmpl, templateError := template.ParseFiles("templates/fetcher.html")
	if templateError != nil {
		log.Printf("Failed to parse template in handlers.go ImageRequest,\n\t\terror: %s", templateError)
		fmt.Fprintf(w, "And unexpected error has occured. Please try again.")
		return
	}
	tmpl.Execute(w, tmpl)
}

// HTTP GET
func DataViewer(w http.ResponseWriter, r *http.Request) {
	tmpl, templateError := template.ParseFiles("templates/viewer.html")
	if templateError != nil {
		log.Printf("Failed to parse template in handlers.go DataViewer,\n\t\terror: %s", templateError)
		fmt.Fprintf(w, "And unexpected error has occured. Please try again.")
		return
	}
	tmpl.Execute(w, tmpl)
}

// HTTP POST
func ImageDataRetrieval(w http.ResponseWriter, r *http.Request) {
	tmpl, templateError := template.ParseFiles("viewer/.html")
	if templateError != nil {
		log.Printf("Failed to parse template in handlers.go ImageDataRetrieval,\n\t\terror: %s", templateError)
		fmt.Fprintf(w, "And unexpected error has occured. Please try again.")
		return
	}
	tmpl.Execute(w, tmpl)
}
