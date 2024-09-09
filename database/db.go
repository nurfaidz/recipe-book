package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"recipebook/models"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	// load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbPort := os.Getenv("PGPORT")
	dbName := os.Getenv("PGDATABASE")
	sslMode := os.Getenv("PGSSLMODE")

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbName, dbPort, sslMode)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	fmt.Println("successfully connected to database")
	db.Debug().AutoMigrate(models.Users{}, models.Recipes{}, models.Comments{}, models.Follows{}, models.Likes{})
}

func GetDB() *gorm.DB {
	return db
}
