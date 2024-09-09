package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"recipebook/database"
	"recipebook/router"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("No file .env found, using system environment variables")
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	database.StartDB()
	router.StartServer().Run(":" + PORT)
}
