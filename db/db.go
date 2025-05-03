package db

import (
	"fmt"
	"log"
	"os"

	"github.com/kithmina1999/eldraread-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database", err)
		panic("Database connection error")
	}
	DB = database
	//auto migrate models
	err = DB.AutoMigrate(
	    &models.User{},
	    &models.Genre{},
	    &models.Novel{},
	    &models.Author{},
	    &models.Tags{},
	    &models.Review{},
	    &models.Comment{},

  	)
	if err != nil {
		fmt.Println("Failed to migrate database", err)
		panic("Database migration error")
	}

}
