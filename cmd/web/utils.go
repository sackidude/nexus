package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetSQLFormattedDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func CalculateVolume(db *sql.DB, pxHeight float64, imageId int64) (float64, error) {
	zeroheight, ml_per_pixel, err := GetTrialPxInformation(db, imageId)
	if err != nil {
		log.Printf("Failed to get trial px information in calculateVolume")
		return 0, errors.New("failed to get trialPxInformation")
	}

	return (pxHeight - float64(zeroheight)) * ml_per_pixel, nil
}

func ExtractInformation(r *http.Request) (id int64, pxHeight float64, err error) {
	id, err1 := strconv.ParseInt(r.Header.Get("id"), 10, 64)
	if err1 != nil {
		log.Printf("Failed to parse trialId from request in ExtractInformation\n\t\terror:%s", err1)
		err := errors.New("failed to parse")
		return 0, 0, err
	}
	pxHeight, err2 := strconv.ParseFloat(r.Header.Get("pxHeight"), 64)
	if err2 != nil {
		log.Printf("Failed to parse pxHeight from request in ExtractInformation\n\t\terror:%s", err2)
		err := errors.New("failed to parse")
		return 0, 0, err
	}
	return id, pxHeight, nil
}
