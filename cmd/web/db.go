package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
