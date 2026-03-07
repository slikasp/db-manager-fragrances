package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/cards"
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
	fmt.Println("---Application started---")

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
	fmt.Println("--Setup complete--")

	// Go throug all IDs after the last found card and look for existing new cards
	fmt.Println("-Looking for new cards-")
	err = cards.FindNewCards(stt, 100)
	if err != nil {
		log.Fatalf("Failed finding new cards: %v", err)
	}

	// Go through newly found cards, try decoding them and add new fragrance entries
	fmt.Println("-Adding missing fragrances-")
	err = fragrances.AddMissingFragrances(stt)
	if err != nil {
		log.Fatalf("Failed adding new fragrances: %v", err)
	}

	// Go through all IDs that have no cards and update if they are now present
	// TODO: make this run in parallel to everything else and remove loggin to file
	fmt.Println("-Checking missing cards started-")
	err = cards.CheckMissingCards(stt)
	if err != nil {
		log.Fatalf("Failed getting new cards: %v", err)
	}

	log.Println("---Application closed---")
}
