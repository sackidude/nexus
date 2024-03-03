package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"
)

func GetNewImageData(db *sql.DB) map[string]string {
	var (
		id          int
		trial       int
		filename    string
		pictureTime string
	)
	rows, sqlError := db.Query("SELECT id, trial, filename, time FROM images ORDER BY request_date ASC LIMIT 1")
	if sqlError != nil {
		log.Printf("Error while querying data base in GetNewImageData\n\t\t%s", sqlError)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &trial, &filename, &pictureTime)
		if err != nil {
			log.Printf("Error scanning lines from database in GetNewImageData\n\t\t%s", err)
		}
	}

	// Update the entry in the db to reflect that someone has seen it
	go func() {
		formattedTime := GetSQLFormattedDateTime()
		query := fmt.Sprintf(`UPDATE Images SET request_date='%s', state='I' WHERE id=%d;`, formattedTime, id)
		_, execError := db.Exec(query)
		if execError != nil {
			log.Printf("Failed to update images, error: %s", execError)
		}
	}()

	path := fmt.Sprintf("trial-%d/%s", trial, filename)

	return map[string]string{
		"id":    fmt.Sprintf("%d", id),
		"trial": fmt.Sprintf("%d", trial),
		"image": filename,
		"time":  pictureTime,
		"path":  path,
	}
}

func SetImageData(db *sql.DB, volume float64, id int64) {
	formattedTime := GetSQLFormattedDateTime()
	query := fmt.Sprintf(
		`UPDATE Images 
		SET volume=(1/(analyzed+1))*(analyzed*volume+%g), analyzed=analyzed+1, request_date='%s', state='D'
		WHERE id=%d;`, volume, formattedTime, id)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to update images in SetImageData. Error: %s", err)
	}
}

func GetTrialPxInformation(db *sql.DB, imageId int64) (zeroheight int64, ml_per_pixel float64, err error) {
	// First get the trial num from
	imageQuery := fmt.Sprintf("SELECT trial FROM Images WHERE id=%d", imageId)
	imageRows, imageQueryError := db.Query(imageQuery)
	if imageQueryError != nil {
		log.Printf("Images query failed in GetTrialPxInformation, error: %s", imageQueryError)
		return 0, 0, errors.New("failed to query images")
	}
	var trialNum int
	for imageRows.Next() {
		rowErr := imageRows.Scan(&trialNum)
		if rowErr != nil {
			log.Printf("Failed to scan row in images Query")
			return 0, 0, errors.New("failed to scan image query")
		}
	}
	rowCloseError := imageRows.Close()
	if rowCloseError != nil {
		log.Printf("Failed to close images rows")
	}

	// Secondly, get the information from the trial
	trialQuery := fmt.Sprintf("SELECT zero_height, ml_per_pixel FROM Trials WHERE trial_num=%d", trialNum)
	trialRows, trialQueryError := db.Query(trialQuery)
	if trialQueryError != nil {
		log.Printf("Failed trial query in GetTrialPxInformation, error: %s", trialQueryError)
		return 0, 0, errors.New("failed to query trials")
	}
	for trialRows.Next() {
		rowErr := trialRows.Scan(&zeroheight, &ml_per_pixel)
		if rowErr != nil {
			log.Printf("Failed to scan row in trial query")
			return 0, 0, errors.New("failed to scan trial query")
		}
	}
	trialRowCloseError := trialRows.Close()
	if trialRowCloseError != nil {
		log.Printf("Failed to close trial rows")
	}

	return zeroheight, ml_per_pixel, nil
}

type Trial struct {
	Done       []Point `json:"done"`
	InProgress []Point `json:"inProgress"`
	Unlabeled  []Point `json:"unlabeled"`
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// TODO: add better error handling to this function
func GetChartData(db *sql.DB) template.JS {
	data := make(map[int]Trial)

	imageRows, imageError := db.Query("SELECT trial, time, volume, state, filename FROM Images")
	if imageError != nil {
		log.Printf("Failed to query images, error: %s", imageError)
		return ""
	}
	var t0 time.Time
	for imageRows.Next() {
		var trialNum int
		var timestring string
		var volume float32
		var state string
		var filename string
		err := imageRows.Scan(&trialNum, &timestring, &volume, &state, &filename)
		if err != nil {
			log.Printf("Failed to scan row, error: %s", err)
			return ""
		}
		timePoint, timeParseErr := SQLTimestampToDatetime(timestring)
		if timeParseErr != nil {
			log.Printf("Failed to parse time")
			return ""
		}
		// Generate the point
		// Just trust that this has run hehe...
		if filename == "1.jpg" {
			t0 = timePoint
		}
		const secondsInHour = 3600.0
		t := float32(timePoint.Sub(t0).Seconds() / secondsInHour)
		point := Point{X: t, Y: volume}

		trial, ok := data[trialNum]
		if !ok {
			data[trialNum] = Trial{}
		}

		switch state {
		case "D":
			trial.Done = append(trial.Done, point)
		case "I":
			trial.InProgress = append(trial.InProgress, point)
		case "U":
			trial.Unlabeled = append(trial.Unlabeled, point)
		default:
			log.Printf("Invalid state in GetChartData")
		}
		data[trialNum] = trial
	}
	closeErr2 := imageRows.Close()
	if closeErr2 != nil {
		log.Printf("Failed to close query error: %s", closeErr2)
	}

	// The returned data should be in the following format
	resultByte, jsonMarshalErr := json.Marshal(data)
	if jsonMarshalErr != nil {
		log.Printf("Failed to marshal item, error: %s", jsonMarshalErr)
		return template.JS("")
	}
	return template.JS(string(resultByte))
}

func GetDBInfo(db *sql.DB) int {
	rows, err := db.Query("SELECT COUNT(*) FROM Images")
	if err != nil {
		log.Printf("Failed to get database infomation, error: %s", err)
	}
	defer rows.Close()
	var entries int
	for rows.Next() {
		err = rows.Scan(&entries)
		if err != nil {
			log.Printf("Failed to scan rows, error: %s", err)
		}
	}
	return entries
}

type viewerTemplate struct {
	Data       template.JS
	Trial_info []int
}

func GetTrialTemplate(db *sql.DB) (viewerTemplate, error) {
	rows, err := db.Query("SELECT trial_num FROM Trials")
	if err != nil {
		return viewerTemplate{}, err
	}
	var res viewerTemplate
	for rows.Next() {
		var temp int
		err = rows.Scan(&temp)
		if err != nil {
			log.Printf("Failed to scan row, error. %s", err)
		}
		res.Trial_info = append(res.Trial_info, temp)
	}

	res.Data = GetChartData(db)

	return res, nil
}
