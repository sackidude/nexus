package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// HTTP GET
func ImageRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, templateError := template.ParseFiles("templates/fetcher.html")
	if templateError != nil {
		log.Printf("Failed to parse template in handlers.go ImageRequest,\n\t\terror: %s", templateError)
		fmt.Fprintf(w, "And unexpected error has occured. Please try again.")
		return
	}

	tmpl.Execute(w, GetNewImageData(db))
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
func ImageDataRetrieval(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Make the query to db
	id, pxHeight, headerError := ExtractInformation(r)
	if headerError != nil {
		log.Printf("Failed to extract information from header in ImageDataRetrieval, error: %s", headerError)
		return
	}
	volume, volumeCalcError := CalculateVolume(db, pxHeight, id)
	if volumeCalcError != nil {
		log.Printf("Failed to calculater volume, error: %s", volumeCalcError)
		return
	}
	go SetImageData(db, volume, id)
	ImageRequest(w, r, db)
}
