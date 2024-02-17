package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: replace all _ with real errorhandling and fallback options
func getNewImageData(db *sql.DB) map[string]string {
	var (
		id       int
		trial    int
		filename string
		time     string
	)
	rows, _ := db.Query("SELECT id, trial, filename, time FROM images ORDER BY last_analyzed ASC LIMIT 1")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &trial, &filename, &time)
		if err != nil {
			println("error reading lines")
		}
	}

	path := fmt.Sprintf("images/trial-%d/%s", trial, filename)

	return map[string]string{
		"id":    fmt.Sprintf("%d", id),
		"trial": fmt.Sprintf("%d", trial),
		"image": filename,
		"time":  time,
		"path":  path,
	}
}

func updateVolumeData(db *sql.DB, id int, volume float64) {
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(
		`UPDATE Images
		SET volume=(1/(analyzed+1))*(analyzed*volume+%g), analyzed=analyzed+1, last_analyzed='%s'
		WHERE id=%d;`, volume, formattedTime, id)

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func calculateVolume(db *sql.DB, trialId int64, pxHeight float64) float64 {
	// get zero height and 1000ml height from db
	var (
		zero_height  int
		ml_per_pixel float64
	)
	query := fmt.Sprintf("SELECT zero_height, ml_per_pixel FROM Trials WHERE trial_num=%d LIMIT 1", trialId)
	rows, _ := db.Query(query)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&zero_height, &ml_per_pixel)
		if err != nil {
			println("error reading lines")
		}
	}

	return (pxHeight - float64(zero_height)) * ml_per_pixel
}

// This is very temporary
func getGraphData() map[string]string {
	return map[string]string{
		"trial-3": "This is where the data would be",
		"trial-4": "This is where the data might be",
	}
}

func main() {
	// Create connection to mysql db
	// Fetch secrets
	secretsJSON, _ := os.ReadFile("../secrets.json")
	type Secrets struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var secrets Secrets
	json.Unmarshal(secretsJSON, &secrets)

	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/nexus", secrets.Username, secrets.Password))
	defer db.Close()

	// Basic static hosting
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// HTMX triggers
	// Preload fragments
	fetcherHTML, _ := os.ReadFile("fragments/fetcher/main.html")
	viewerHTML, _ := os.ReadFile("fragments/viewer/main.html")

	// HTTP Requests
	imageRequest := func(w http.ResponseWriter, r *http.Request) {
		// Request the relevant information to db.
		// Template html fragment with the information from db
		tmpl, _ := template.New("t").Parse(string(fetcherHTML))
		imageData := getNewImageData(db)
		tmpl.Execute(w, imageData)
	}
	http.HandleFunc("/image-request", imageRequest)

	dataViewer := func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("t").Parse(string(viewerHTML))
		graphData := getGraphData()
		tmpl.Execute(w, graphData)
	}
	http.HandleFunc("/data-view", dataViewer)

	// HTTP Posts
	imageDataRetrival := func(w http.ResponseWriter, r *http.Request) {
		trialId, _ := strconv.ParseInt(r.Header.Get("trialId"), 10, 32)
		id, _ := strconv.ParseInt(r.Header.Get("id"), 10, 32)
		id32 := int(id)
		pxHeightString := r.Header.Get("pxHeight")
		pxHeight, _ := strconv.ParseFloat(pxHeightString, 64)

		volume := calculateVolume(db, trialId, pxHeight)
		fmt.Printf("User image data received with volume: %.1f\n", volume)
		updateVolumeData(db, id32, volume)

		imageData := getNewImageData(db)
		tmpl, _ := template.New("t").Parse(string(fetcherHTML))
		tmpl.Execute(w, imageData)
	}
	http.HandleFunc("/user-image-data", imageDataRetrival)

	http.ListenAndServe("localhost:8080", nil)
}
