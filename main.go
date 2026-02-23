package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

func main() {
	// logFile, err := setupLogging("app.log")
	// if err != nil {
	// 	log.Fatalf("failed to setup logging: %v", err)
	// }
	// defer logFile.Close()

	maxFragrances := int32(123330)

	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Load the database
	dbtx, err := sql.Open("postgres", cfg.RemoteDbURL)
	dbQueries := database.New(dbtx)

	// Create database struct to be passed to functions
	stt := &config.State{
		DB:        dbQueries,
		CurrentID: 1,
	}

	// start from 0 to maxFragrances
	//make function of the below
	// try downloading one image
	// update db
}
