package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func Connect() {
	// Database connection string
	dsn := "host=localhost user=postgres dbname=bookAPI sslmode=disable password=postgress port=5432"

	// Opening connection to database
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	return db
}
