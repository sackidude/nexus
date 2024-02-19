package main

import (
	"log"
	"net/http"
	"os"

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

	log.Printf("Login details:\n\tUsername: %s \n\tPassword: %s", DBCredentials.Username, DBCredentials.Password)

	// Static file hosting.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	listeningErr := http.ListenAndServe("localhost:8080", nil)
	if listeningErr != nil {
		log.Printf("Error in serving and listening:\n\t\terror: %s", listeningErr)
	}
}
