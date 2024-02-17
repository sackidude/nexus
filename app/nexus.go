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
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: replace all _ with real errorhandling and fallback options
func getNewImageData(db *sql.DB) map[string]string {
	var (
		id          int
		trial       int
		filename    string
		pictureTime string
	)
	rows, _ := db.Query("SELECT id, trial, filename, time FROM images ORDER BY last_seen ASC LIMIT 1")
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &trial, &filename, &pictureTime)
		if err != nil {
			println("error reading lines")
		}
	}

	// Update the entry in the db to reflect that someone has seen it
	go func() {
		formattedTime := time.Now().Format("2006-01-02 15:04:05")
		query := fmt.Sprintf(`UPDATE Images SET last_seen='%s' WHERE id=%d;`, formattedTime, id)
		db.Exec(query)
	}()

	path := fmt.Sprintf("images/trial-%d/%s", trial, filename)

	return map[string]string{
		"id":    fmt.Sprintf("%d", id),
		"trial": fmt.Sprintf("%d", trial),
		"image": filename,
		"time":  pictureTime,
		"path":  path,
	}
}

func updateVolumeData(db *sql.DB, id int, volume float64) {
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(
		`UPDATE Images
		SET volume=(1/(analyzed+1))*(analyzed*volume+%g), analyzed=analyzed+1, last_seen='%s'
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

func generateGraphJSON(db *sql.DB) template.JS {
	// Request to db
	// First get all the trials
	var trials []int
	{
		rows, _ := db.Query("SELECT trial_num FROM Trials")
		defer rows.Close()

		var temp int
		for rows.Next() {
			rows.Scan(&temp)
			trials = append(trials, temp)
		}
	}

	// Get all of the datapoints in first trial to begin with
	{
		query := fmt.Sprintf("SELECT time, volume FROM Images WHERE trial=%d", trials[0])
		rows, _ := db.Query(query)
		defer rows.Close()

		type datapoint struct {
			datetime time.Time
			volume   float64
		}
		var datapoints []datapoint
		layout := "2006-01-02 15:04:05"
		for rows.Next() {
			var timestr string
			var dp datapoint
			rows.Scan(&timestr, &dp.volume)

			dp.datetime, _ = time.Parse(layout, timestr)
			datapoints = append(datapoints, dp)
		}
		// convert datapoints into the format js format
		// That is [{x:1,y:2},{x:2,y:3}]
		var jsonStrings []string
		t0 := datapoints[0].datetime.Unix()
		for _, v := range datapoints {
			secondsInHour := 3600.0
			hours := float64(v.datetime.Unix()-t0) / secondsInHour
			tempStr := fmt.Sprintf("{x:%.2f,y:%.1f}", hours, v.volume)
			jsonStrings = append(jsonStrings, tempStr)
		}
		result := fmt.Sprintf(`[%s]`, strings.Join(jsonStrings, ","))
		return template.JS(result)
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
		graphData := generateGraphJSON(db)
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
