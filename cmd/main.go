package main

import (
	"arxivhub/internal/app"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	application, err := app.NewApp()
	if err != nil {
		log.Fatalln(err.Error())
	}

	application.Run(":8080")
}
