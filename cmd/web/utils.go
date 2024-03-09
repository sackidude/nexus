package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetSQLFormattedDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func SQLTimestampToDatetime(timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, timeStr)
}

func CalculateVolume(db *sql.DB, pxHeight float64, imageId int64) (float64, error) {
	zeroheight, ml_per_pixel, err := GetTrialPxInformation(db, imageId)
	if err != nil {
		return 0, fmt.Errorf("GetTrialPxInformation: %w", err)
	}

	return (pxHeight - float64(zeroheight)) * ml_per_pixel, nil
}

func ExtractInformation(r *http.Request) (id int64, pxHeight float64, err error) {
	id, err = strconv.ParseInt(r.Header.Get("id"), 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}
	pxHeight, err = strconv.ParseFloat(r.Header.Get("pxHeight"), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("strconv.ParseFloat: %w", err)
	}
	return id, pxHeight, nil
}

func GetTrialNumFromHeader(r *http.Request) (int, error) {
	trial_num, err := strconv.ParseInt(r.Header.Get("trial_num"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt(r.Header.Get('trial_num'), 10, 32): %w", err)
	}
	return int(trial_num), nil
}
