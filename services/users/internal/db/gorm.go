package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}
