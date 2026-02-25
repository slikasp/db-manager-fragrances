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
	// TODO get this from fragrantica main page
	maxFragrances := int32(124000)

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

	// Runs some kind of function with a loop
	// TODO: need single maintenance function when fragrances DB is up to date
	err = cards.DownloadAllCards(stt)

	fmt.Println(err)

	// doesn't work if you CTRL+C out of the loop, manually update config.json in that case
	cfg.CurrentID = stt.CurrentID
	err = config.Write(cfg)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error writing config: %v", err)
	}
}
