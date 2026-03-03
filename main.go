package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
	"github.com/slikasp/dbmanfrags/fragrances"
)

func setup() (*config.State, error) {
	// Read config
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	// Load the database
	dbtx, err := sql.Open("postgres", cfg.RemoteDbURL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(dbtx)

	// Create database struct to be passed to functions
	state := &config.State{
		DB: dbQueries,
	}

	// Set ID of last card from the database
	state.LastID, err = state.DB.GetLastCardID(context.Background())
	if err != nil {
		return nil, err
	}

	return state, nil
}

func main() {
	// Log
	logFile, err := setupLogging("app.log")
	if err != nil {
		log.Fatalf("failed to setup logging: %v", err)
	}
	defer logFile.Close()

	// Setup
	stt, err := setup()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	log.Println("application started")

	// Run
	err = fragrances.AddMissingFragrances(stt)
	if err != nil {
		log.Fatalf("Failed running commands: %v", err)
	}
}
