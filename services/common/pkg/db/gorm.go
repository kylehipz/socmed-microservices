package db

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(log *zap.Logger, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	// Get the underlying *sql.DB object from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Ping the database to force a connection and check if it's alive.
	// This will block until the connection is established or fails.
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// --- (Optional) Configure Connection Pool ---
	// These are good defaults for a production-ready app.
	sqlDB.SetMaxIdleConns(10)           // Set max idle connections
	sqlDB.SetMaxOpenConns(100)          // Set max open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set max connection lifetime

	log.Info("Database connection established successfully.")
	return db, nil
}

