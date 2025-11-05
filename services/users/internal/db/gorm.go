package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Get the underlying *sql.DB object from GORM
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// Ping the database to force a connection and check if it's alive.
	// This will block until the connection is established or fails.
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// --- (Optional) Configure Connection Pool ---
	// These are good defaults for a production-ready app.
	sqlDB.SetMaxIdleConns(10)           // Set max idle connections
	sqlDB.SetMaxOpenConns(100)          // Set max open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set max connection lifetime

	log.Println("Database connection established successfully.")
	return db
}
