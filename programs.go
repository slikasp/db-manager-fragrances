package main

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
	"github.com/slikasp/dbmanfrags/cards"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/fragrances"
)

// Manual update of DB items only
func ManualDbUpdate(frags *config.Frags) {
	checkMissingCards(frags)
	lookForNewCards(frags)
	addMissingFragrances(frags)
}

// Service to keep updating DB with details from the web, maxRequests number of requests per hour
func ScraperService(frags *config.Frags, maxRequests int) {
	c := cron.New()

	c.AddFunc("30 8-23 * * *", func() {
		updatePerfumers(frags, maxRequests)
	})

	c.Start()

	// Keep program alive
	select {}

	//TODO:
	// 1 - get all perfumers updated (~1000)
	// 2 - get all existing perfumes updated (~60000)
	// 3 - make a function to keep updating the oldes fragrances in the database
	// 4 - combine with ManualDbUpdate
}

// Go through all IDs that have no cards and update if they are now present
// This part takes quite a long time so I would run this rarely and only report on the progress
// >40k missing cards, takes almost 3 hours now, yikes!
func checkMissingCards(frags *config.Frags) {
	// TODO: make this run in parallel to everything else and remove loggin to file
	log.Println("> Checking missing cards - start >")
	fmt.Println("-Checking missing cards-")

	err := cards.CheckMissingCards(frags)
	if err != nil {
		log.Fatalf("Failed getting new cards: %v", err)
	}
	log.Println("< Checking missing cards - end <")
}

// Go throug all IDs after the last found card and look for new cards
func lookForNewCards(frags *config.Frags) {
	// With low number of cards to check you can hit a patch of no cards that will prevent you form fiding new ones
	// 100 - too low
	// 1000 - a bit too much
	// 300 - a safe option, takes less than 2 minutes if none are found
	log.Println("> Looking for new cards - start >")
	fmt.Println("-Looking for new cards-")

	err := cards.FindNewCards(frags, 300)
	if err != nil {
		log.Fatalf("Failed finding new cards: %v", err)
	}
	log.Println("< Looking for new cards - end <")
}

// Go through newly found cards, try decoding them and add new fragrance entries
func addMissingFragrances(frags *config.Frags) {
	log.Println("> Adding missing fragrances - start >")
	fmt.Println("-Adding missing fragrances-")

	err := fragrances.AddMissingFragrances(frags)
	if err != nil {
		log.Fatalf("Failed adding new fragrances: %v", err)
	}

	log.Println("< Adding missing fragrances - end <")
}

// Go through new fragrances, find relevant data and update them
func addMissingDetails(frags *config.Frags, numRequests int) {
	log.Println("> Adding missing details - start >")
	fmt.Println("-Adding missing details-")

	// delay the start of the program 1-5 minutes
	fragrances.SpamDelay(60, 300)

	err := fragrances.UpdateFragrances(frags, numRequests)
	if err != nil {
		log.Fatalf("Failed adding details: %v", err)
	}

	log.Println("< Adding missing details - end <")
}

// do not use after initial list of perfumers is complete - new perfumer to be added together with fragrance
// TODO: remake this so all perfumers need to be updated - add 'updated' column
func updatePerfumers(frags *config.Frags, numRequests int) {
	log.Println("> Updating perfumers - start >")
	fmt.Println("-Updating perfumers-")

	// delay the start of the program 1-5 minutes
	fragrances.SpamDelay(60, 300)

	err := fragrances.UpdatePerfumers(frags, numRequests)
	if err != nil {
		log.Fatalf("Failed updating perfumers: %v", err)
	}
	log.Println("< Updating perfumers - end <")
}
