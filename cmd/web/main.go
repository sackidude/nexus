package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Started!")
	err := godotenv.Load("../../config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBCredentials := struct {
		Username string
		Password string
	}{os.Getenv("SQL_USERNAME"), os.Getenv("SQL_PASSWORD")}

	log.Printf("Login details:\n\tUsername: %s \n\tPassword:%s", DBCredentials.Username, DBCredentials.Password)
}
