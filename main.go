package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/cards"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

func main() {
	logFile, err := setupLogging("app.log")
	if err != nil {
		log.Fatalf("failed to setup logging: %v", err)
	}
	defer logFile.Close()

	// TODO get this from fragrantica main page
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
		CurrentID: cfg.CurrentID,
		LastID:    maxFragrances,
	}

	l := cards.GetAllCards(*stt)

	fmt.Println(l)

	cfg.CurrentID = stt.CurrentID
	err = config.Write(cfg)
	if err != nil {
		log.Fatalf("error writing config: %v", err)
	}
}
