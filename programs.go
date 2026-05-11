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
func ManualDbUpdate(db *config.Database) {
	checkMissingCards(db)
	lookForNewCards(db)
	addMissingFragrances(db)
}

// Service to keep updating DB with details from the web, maxRequests number of requests per hour
func ScraperService(db *config.Database, maxRequests int) {
	c := cron.New()

	schedule := "30 8-23 * * *"

	// only run these in the production environment where local card cache is located
	if db.BuildEnv == "prod" {
		// this will run once a week to check if we have downloaded all cards listed in DB and if any missing cards were uploaded to the web
		c.AddFunc("0 0 * * 1", func() {
			checkExistingCards(db)
			checkMissingCards(db)
		})

		// this will once every morning to look for new fragrances
		c.AddFunc("0 7 * * *", func() {
			lookForNewCards(db)
			addMissingFragrances(db)
		})

		schedule = "00 8-23 * * *"
	}

	// this will run every waking hour every day and keep updating maxRequests*16 fragrance items every day
	c.AddFunc(schedule, func() {
		fragrances.SpamDelay(60, 300)
		updateFragranceDetails(db, maxRequests)
	})

	c.Start()

	// Keep program alive
	select {}

	//TODO:
	// 1 - get all perfumers updated (~1000) - done
	// 2 - get all existing perfumes updated (~60000) - in progress
	// 3 - make a function to keep updating the oldest fragrances in the database - done
	// 4 - combine with ManualDbUpdate - done
}

// Go through all IDs that have no cards and update if they are now present
// This part takes quite a long time so I would run this rarely and only report on the progress
// >40k missing cards, takes almost 3 hours now, yikes!
func checkMissingCards(db *config.Database) {
	// TODO: make this run in parallel to everything else and remove loggin to file
	log.Println("> Checking missing cards - start >")
	fmt.Println("-Checking missing cards-")

	err := cards.CheckMissingCards(db)
	if err != nil {
		log.Fatalf("Failed getting new cards: %v", err)
	}
	log.Println("< Checking missing cards - end <")
}

func checkExistingCards(db *config.Database) {
	// TODO: make this run in parallel to everything else and remove loggin to file
	log.Println("> Checking downloaded cards - start >")
	fmt.Println("-Checking downloaded cards-")

	err := cards.CheckExistingCards(db)
	if err != nil {
		log.Fatalf("Failed checking downloaded cards: %v", err)
	}
	log.Println("< Checking downloaded cards - end <")
}

// Go throug all IDs after the last found card and look for new cards
func lookForNewCards(db *config.Database) {
	// With low number of cards to check you can hit a patch of no cards that will prevent you form fiding new ones
	// 100 - too low
	// 1000 - a bit too much
	// 300 - a safe option, takes less than 2 minutes if none are found
	log.Println("> Looking for new cards - start >")
	fmt.Println("-Looking for new cards-")

	err := cards.FindNewCards(db, 300)
	if err != nil {
		log.Fatalf("Failed finding new cards: %v", err)
	}
	log.Println("< Looking for new cards - end <")
}

// Go through newly found cards, try decoding them and add new fragrance entries
func addMissingFragrances(db *config.Database) {
	log.Println("> Adding missing fragrances - start >")
	fmt.Println("-Adding missing fragrances-")

	err := fragrances.AddMissingFragrances(db)
	if err != nil {
		log.Fatalf("Failed adding new fragrances: %v", err)
	}

	log.Println("< Adding missing fragrances - end <")
}

// Go through new fragrances, find relevant data and update them
func updateFragranceDetails(db *config.Database, numRequests int) {
	log.Println("> Adding missing details - start >")
	fmt.Println("-Adding missing details-")

	err := fragrances.UpdateFragrances(db, numRequests)
	if err != nil {
		log.Fatalf("UpdateFragrances: %v", err)
	}

	log.Println("< Adding missing details - end <")
}

// do not use after initial list of perfumers is complete - new perfumer to be added together with fragrance
// TODO: remake this so all perfumers need to be updated - add 'updated' column
func updatePerfumers(db *config.Database, numRequests int) {
	log.Println("> Updating perfumers - start >")
	fmt.Println("-Updating perfumers-")

	err := fragrances.UpdatePerfumers(db, numRequests)
	if err != nil {
		log.Fatalf("Failed updating perfumers: %v", err)
	}
	log.Println("< Updating perfumers - end <")
}
