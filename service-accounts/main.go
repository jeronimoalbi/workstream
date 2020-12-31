package main

import (
	"fmt"
	"os"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database connection
var database *gorm.DB

// KUSANAGI service
var service = kusanagi.NewService()

// DSN creates a new data source name to connect to the database.
func DSN() string {
	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		user = "workstream"
	}

	password := os.Getenv("DATABASE_PASWORD")
	if password == "" {
		password = "workstream"
	}

	name := os.Getenv("DATABASE_NAME")
	if name == "" {
		name = "workstream"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, name)
}

func main() {
	// Connect to the database
	dsn := DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
	} else {
		log.Debugf("Connected to database: %s", dsn)
		database = db
	}

	// Run the KUSANAGI service
	service.Run()
}
