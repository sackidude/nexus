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

func GetNewImageData(db *sql.DB) (map[string]string, error) {
	var (
		id          int
		trial       int
		filename    string
		pictureTime string
	)
	rows, err := db.Query("SELECT id, trial, filename, time FROM images ORDER BY request_date LIMIT 1")
	if err != nil {
		return map[string]string{}, fmt.Errorf("db.Query: %w", err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &trial, &filename, &pictureTime)
		if err != nil {
			return map[string]string{}, fmt.Errorf("rows.Scan: %w", err)
		}
	}

	// Update the entry in the db to reflect that someone has seen it
	go func() {
		formattedTime := GetSQLFormattedDateTime()
		query := fmt.Sprintf(`UPDATE Images SET request_date='%s', state='I' WHERE id=%d;`, formattedTime, id)
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("GetNewImageData: Anon function: db.Exec: %s", err)
		}
	}()

	path := fmt.Sprintf("trial-%d/%s", trial, filename)

	return map[string]string{
		"id":    fmt.Sprintf("%d", id),
		"trial": fmt.Sprintf("%d", trial),
		"image": filename,
		"time":  pictureTime,
		"path":  path,
	}, nil
}

func SetImageData(db *sql.DB, volume float64, id int64) error {
	formattedTime := GetSQLFormattedDateTime()
	query := fmt.Sprintf(
		`UPDATE Images 
		SET volume=(1/(analyzed+1))*(analyzed*volume+%g), analyzed=analyzed+1, request_date='%s', state='D'
		WHERE id=%d;`, volume, formattedTime, id)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func GetTrialPxInformation(db *sql.DB, imageId int64) (zeroheight int64, ml_per_pixel float64, err error) {
	// First get the trial num from
	imageQuery := fmt.Sprintf("SELECT trial FROM Images WHERE id=%d", imageId)
	imageRows, err := db.Query(imageQuery)
	if err != nil {
		return 0, 0, fmt.Errorf("db.Query(imageQuery): %w", err)
	}
	var trialNum int
	for imageRows.Next() {
		err = imageRows.Scan(&trialNum)
		if err != nil {
			return 0, 0, fmt.Errorf("imageRows.Scan: %w", err)
		}
	}

	// Secondly, get the information from the trial
	trialQuery := fmt.Sprintf("SELECT zero_height, ml_per_pixel FROM Trials WHERE trial_num=%d", trialNum)
	trialRows, err := db.Query(trialQuery)
	if err != nil {
		return 0, 0, fmt.Errorf("db.Query(trialQuery): %w", err)
	}
	for trialRows.Next() {
		err := trialRows.Scan(&zeroheight, &ml_per_pixel)
		if err != nil {
			return 0, 0, fmt.Errorf("trialRows.Scan: %w", err)
		}
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
func GetChartData(db *sql.DB) (template.JS, error) {
	data := make(map[int]Trial)

	imageRows, err := db.Query("SELECT trial, time, volume, state, filename FROM Images")
	if err != nil {
		return "", fmt.Errorf("db.Query: %w", err)
	}
	var t0 time.Time
	for imageRows.Next() {
		var trialNum int
		var timestring string
		var volume float32
		var state string
		var filename string
		err = imageRows.Scan(&trialNum, &timestring, &volume, &state, &filename)
		if err != nil {
			return "", fmt.Errorf("imageRows.Scan: %w", err)
		}
		timePoint, err := SQLTimestampToDatetime(timestring)
		if err != nil {
			return "", fmt.Errorf("SQLTimestampToDatetime: %w", err)
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
			return "", errors.New("invalid trial state from db")
		}
		data[trialNum] = trial
	}

	// The returned data should be in the following format
	resultByte, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}
	return template.JS(string(resultByte)), nil
}

func GetDBInfo(db *sql.DB) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM Images")
	if err != nil {
		return 0, fmt.Errorf("db.Query: %w", err)
	}
	var entries int
	for rows.Next() {
		err = rows.Scan(&entries)
		if err != nil {
			return 0, fmt.Errorf("rows.Scan: %w", err)
		}
	}
	return entries, nil
}

type viewerTemplate struct {
	Data       template.JS
	Trial_info []int
}

func GetTrialTemplate(db *sql.DB) (viewerTemplate, error) {
	rows, err := db.Query("SELECT trial_num FROM Trials")
	if err != nil {
		return viewerTemplate{}, fmt.Errorf("db.Query: %w", err)
	}
	var res viewerTemplate
	for rows.Next() {
		var temp int
		err = rows.Scan(&temp)
		if err != nil {
			return viewerTemplate{}, fmt.Errorf("rows.Scan: %w", err)
		}
		res.Trial_info = append(res.Trial_info, temp)
	}

	res.Data, err = GetChartData(db)
	if err != nil {
		return viewerTemplate{}, fmt.Errorf("GetChartData: %w", err)
	}

	return res, nil
}

type trialTemplate struct {
	TrialNum       int
	StartTimestamp string
	YeastAmount    int
	SugarAmount    int
	Stirring       bool
}

func GetFullTrialInfo(db *sql.DB, trial_num int) (trialTemplate, error) {
	// Query the trial for info.
	trialRows, err := db.Query("SELECT start, yeast_amount, sugar_amount, stirring FROM Trials LIMIT 1")
	if err != nil {
		return trialTemplate{}, fmt.Errorf("db.Query: %w", err)
	}
	var start string
	var yeast_amount int
	var sugar_amount int
	var stirringInt int
	for trialRows.Next() {
		err = trialRows.Scan(&start, &yeast_amount, &sugar_amount, &stirringInt)
		if err != nil {
			return trialTemplate{}, fmt.Errorf("trialRows.Scan: %w", err)
		}
	}
	stirring := stirringInt != 0
	return trialTemplate{trial_num, start, yeast_amount, sugar_amount, stirring}, nil
}
