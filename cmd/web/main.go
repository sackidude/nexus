package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Started!")

	// Dotenv loading and getting of credentials.
	err := godotenv.Load("../../config/.env")
	if err != nil {
		log.Fatalf("main: godotenv.Load: %s", err)
	}

	DBCredentials := struct {
		Username string
		Password string
	}{os.Getenv("SQL_USERNAME"), os.Getenv("SQL_PASSWORD")}

	// Connect to db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/nexus", DBCredentials.Username, DBCredentials.Password))
	if err != nil {
		log.Fatalf("main: sql.Open: %s", err)
	}
	log.Println("Successfully connected to database")

	// Static file hosting.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// HTMX
	// HTTP GET handling
	http.HandleFunc("/image-request", func(w http.ResponseWriter, r *http.Request) {
		ImageRequest(w, r, db)
	})
	http.HandleFunc("/viewer", func(w http.ResponseWriter, r *http.Request) {
		DataViewer(w, r, db)
	})
	http.HandleFunc("/start-page", func(w http.ResponseWriter, r *http.Request) {
		StartPage(w, r, db)
	})

	// HTTP POST handling
	http.HandleFunc("/user-image-data", func(w http.ResponseWriter, r *http.Request) {
		ImageDataRetrieval(w, r, db)
	})

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Printf("main: http.ListenAndServe: %s", err)
	}
}
