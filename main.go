package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/fragrances"
)

func main() {
	log.Println("---Application started---")
	fmt.Println("---Application started---")

	logFile, err := setupLogging("app.log")
	if err != nil {
		log.Fatalf("failed to setup logging: %v", err)
	}
	defer logFile.Close()

	// Setup
	frags, err := config.Setup()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	log.Println("--Setup complete--")
	fmt.Println("--Setup complete--")

	// TODO: rework the script to go through IDs and checking what needs to be done
	// requires a function that does everything

	// Go through all IDs that have no cards and update if they are now present
	// This part takes quite a long time so I would run this rarely and only report on the progress
	// >40k missing cards, takes almost 3 hours now, yikes!
	// TODO: make this run in parallel to everything else and remove loggin to file
	log.Println("-Checking missing cards-")
	fmt.Println("-Checking missing cards-")
	// err = cards.CheckMissingCards(stt)
	// if err != nil {
	// 	log.Fatalf("Failed getting new cards: %v", err)
	// }

	// Go throug all IDs after the last found card and look for new cards
	// With low number of cards to check you can hit a patch of no cards that will prevent you form fiding new ones
	// 100 - too low
	// 1000 - a bit too much
	// 300 - a safe option, takes less than 2 minutes if none are found
	log.Println("-Looking for new cards-")
	fmt.Println("-Looking for new cards-")
	// err = cards.FindNewCards(frags, 300)
	// if err != nil {
	// 	log.Fatalf("Failed finding new cards: %v", err)
	// }

	// Go through newly found cards, try decoding them and add new fragrance entries
	log.Println("-Adding missing fragrances-")
	fmt.Println("-Adding missing fragrances-")
	// err = fragrances.AddMissingFragrances(frags)
	// if err != nil {
	// 	log.Fatalf("Failed adding new fragrances: %v", err)
	// }

	// Go through new fragrances (which field to check?) and find relevant data in the html
	log.Println("-Adding missing details-")
	fmt.Println("-Adding missing details-")
	err = fragrances.AddMissingDetails(frags)
	if err != nil {
		log.Fatalf("Failed adding details: %v", err)
	}

	// Go through all existing fragrances and check if it needs updating? maybe check the user score

	log.Println("---Application closed---")
	fmt.Println("---Application closed---")
}
