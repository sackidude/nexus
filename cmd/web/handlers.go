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
		fmt.Fprint(w, "An unexpected error has occured. Please try again.")
		return
	}

	imageData, err := GetNewImageData(db)

	if err != nil {
		log.Printf("ImageRequest: GetNewImageData: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again.")
		return
	}

	tmpl.Execute(w, imageData)
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
		fmt.Fprint(w, "An unexpected error has occured. Please try again.")
		return
	}
	tmpl.Execute(w, templateData)
}

// HTTP GET
func StartPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("templates/startpage.html")
	if err != nil {
		log.Printf("StartPage: template.ParseFiles: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again.")
		return
	}
	databaseInfo, err := GetDBInfo(db)
	if err != nil {
		log.Printf("StartPage: GetDBInfo: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again")
		return
	}
	tmpl.Execute(w, databaseInfo)
}

// HTTP POST
func ImageDataRetrieval(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Make the query to db
	id, pxHeight, err := ExtractInformation(r)
	if err != nil {
		log.Printf("ImageDataRetrieval: ExtractInformation: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again")
		return
	}
	volume, err := CalculateVolume(db, pxHeight, id)
	if err != nil {
		log.Printf("ImageDataRetrieval: CalculateVolume: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again")
		return
	}
	go SetImageData(db, volume, id)
	ImageRequest(w, r, db)
}

func DataBaseEntries(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	databaseInfo, err := GetDBInfo(db)
	if err != nil {
		log.Printf("DataBaseEntries: GetDBInfo: %s", err)
		fmt.Fprint(w, "An unexpected error has occured. Please try again")
		return
	}
	fmt.Fprintf(w, "%d", databaseInfo)
}
