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
	envloadingerr := godotenv.Load("../../config/.env")
	if envloadingerr != nil {
		log.Fatalf("Error loading .env file\n\t\t error: %s", envloadingerr)
	}

	DBCredentials := struct {
		Username string
		Password string
	}{os.Getenv("SQL_USERNAME"), os.Getenv("SQL_PASSWORD")}

	// Connect to db
	db, dataBaseErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/nexus", DBCredentials.Username, DBCredentials.Password))
	if dataBaseErr != nil {
		log.Fatalf("Error while connecting to database\n\t\terror: %s", dataBaseErr)
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	// Static file hosting.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// HTMX
	// HTTP GET handling
	http.HandleFunc("/image-request", func(w http.ResponseWriter, r *http.Request) {
		ImageRequest(w, r, db)
	})
	http.HandleFunc("/data-view", DataViewer)

	// HTTP POST handling
	http.HandleFunc("/user-image-data", func(w http.ResponseWriter, r *http.Request) {
		ImageDataRetrieval(w, r, db)
	})

	listeningErr := http.ListenAndServe("localhost:8080", nil)
	if listeningErr != nil {
		log.Printf("Error in serving and listening:\n\t\terror: %s", listeningErr)
	}
}
