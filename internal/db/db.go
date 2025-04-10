package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/G4L1L10/admin-dashboard-backend/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	cfg := config.AppConfig // Access your loaded config

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to the database successfully!")
}
