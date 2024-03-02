package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"strings"
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

// TODO: add better error handling to this function
func GetChartData(db *sql.DB) template.JS {
	// Get all the trials
	var trials []int
	trialRows, trialError := db.Query("SELECT trial_num FROM Trials")
	if trialError != nil {
		log.Printf("Failed to query trial error:%s", trialError)
		return "" // TODO: maybe change this?
	}
	for trialRows.Next() {
		var temp int
		err := trialRows.Scan(&temp)
		if err != nil {
			log.Printf("Error during scan of trials, error: %s", err)
			return ""
		}
		trials = append(trials, temp)
	}
	closeErr := trialRows.Close()
	if closeErr != nil {
		log.Printf("Failed to close query error: %s", closeErr)
	}

	type datapoint struct {
		timePoint time.Time
		volume    float64
		state     string
	}

	trialData := make(map[int][]datapoint)

	imageRows, imageError := db.Query("SELECT trial, time, volume, state FROM Images")
	if imageError != nil {
		log.Printf("Failed to query images, error: %s", imageError)
		return ""
	}
	for imageRows.Next() {
		var trial int
		var timestring string
		var dp datapoint
		err := imageRows.Scan(&trial, &timestring, &dp.volume, &dp.state)
		if err != nil {
			log.Printf("Failed to scan row, error: %s", err)
			return ""
		}
		var timeParseErr error
		dp.timePoint, timeParseErr = SQLTimestampToDatetime(timestring)
		if timeParseErr != nil {
			log.Printf("Failed to parse time")
			return ""
		}
		_, ok := trialData[trial]
		if !ok {
			trialData[trial] = []datapoint{dp}
		} else {
			trialData[trial] = append(trialData[trial], dp)
		}
	}
	closeErr2 := imageRows.Close()
	if closeErr2 != nil {
		log.Printf("Failed to close query error: %s", closeErr2)
	}

	// x = hour y = volume
	//[{x:1, y:2}, {x:2, y:3}]
	var sb strings.Builder
	_, err := sb.WriteString("[")
	if err != nil {
		log.Printf("Failed to push to string builder, error: %s", err)
		return ""
	}

	for trial_num, dpArr := range trialData {
		if trial_num == 3 {
			t0 := trialData[trial_num][0].timePoint
			for i, dp := range dpArr {
				// Timestamp to hours after t0
				secondsInHour := 3600.0
				hour := dp.timePoint.Sub(t0).Seconds() / secondsInHour
				var comma string
				if i != 0 {
					comma = ","
				} else {
					comma = ""
				}
				_, writeStrErr := sb.WriteString(fmt.Sprintf("%s{x:%.3f,y:%.3f}", comma, hour, dp.volume))
				if writeStrErr != nil {
					log.Printf("Failed to push to string builder, error: %s", err)
					return ""
				}
			}
		}

	}
	sb.WriteString("]")
	return template.JS(sb.String())
}
