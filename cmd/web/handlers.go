package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// HTTP GET
func ImageRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/fetcher.html")
	if err != nil {
		log.Printf("ImageRequest: template.ParseFiles: %s", err)
		fmt.Fprint(w, "And unexpected error has occured. Please try again.")
		return
	}

	tmpl.Execute(w, GetNewImageData(db))
}

// HTTP GET
func DataViewer(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/viewer.html")
	if err != nil {
		log.Printf("DataViewer: template.ParseFiles: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again.")
		return
	}
	templateData, err := GetTrialTemplate(db)
	if err != nil {
		log.Printf("DataViewer: GetTrialTemplate: %s", err)
		fmt.Fprintf(w, "An unexpected error has occured. Please try again.")
		return
	}
	tmpl.Execute(w, templateData)
}

// HTTP GET
func StartPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/startpage.html")
	if err != nil {
		log.Printf("StartPage: template.ParseFiles: %s", err)
		fmt.Fprintf(w, "An unexpected error has occured. Please try again.")
		return
	}
	databaseInfo := GetDBInfo(db)
	tmpl.Execute(w, databaseInfo)
}

// HTTP POST
func ImageDataRetrieval(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Make the query to db
	id, pxHeight, err := ExtractInformation(r)
	if err != nil {
		log.Printf("ImageDataRetrieval: ExtractInformation: %s", err)
		return
	}
	volume, err := CalculateVolume(db, pxHeight, id)
	if err != nil {
		log.Printf("ImageDataRetrieval: CalculateVolume: %s", err)
		return
	}
	go SetImageData(db, volume, id)
	ImageRequest(w, r, db)
}
