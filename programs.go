package main

import (
	"github.com/robfig/cron/v3"
	"github.com/slikasp/dbmanfrags/cards"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/fragrances"
)

// Service to keep updating DB with details from the web, maxRequests number of requests per hour
func ScraperService(db *config.Database, numRequests int) {
	c := cron.New()

	schedule := "30 8-23 * * *"

	// only run these in the production environment where local card cache is located
	if db.BuildEnv == "prod" {
		// this will run once a week to check if we have downloaded all cards listed in DB and if any missing cards were uploaded to the web
		c.AddFunc("0 0 * * 1", func() {
			db.Logger.Info("checking downloaded cards")
			cards.CheckExistingCards(db)

			// Go through all IDs that have no cards and update if they are now present
			// This part takes quite a long time so I would run this rarely and only report on the progress
			// >40k missing cards, takes almost 3 hours now, yikes!
			db.Logger.Info("checking missing cards")
			cards.CheckMissingCards(db)
		})

		// this will once every morning to look for new fragrances
		c.AddFunc("0 7 * * *", func() {
			// Go throug all IDs after the last found card and look for new cards
			db.Logger.Info("looking for new cards")
			cards.FindNewCards(db, 300)

			// Go through newly found cards, try decoding them and add new fragrance entries
			db.Logger.Info("adding missing fragrances")
			fragrances.AddMissingFragrances(db)
		})

		schedule = "00 8-23 * * *"
	}

	// this will run every waking hour every day and keep updating maxRequests*16 fragrance items every day
	c.AddFunc(schedule, func() {
		// Go through new fragrances, find relevant data and update them
		db.Logger.Info("adding missing details")
		fragrances.SpamDelay(60, 300)
		fragrances.UpdateFragrances(db, numRequests)
	})

	c.Start()

	// Keep program alive
	select {}
}

// do not use after initial list of perfumers is complete - new perfumer to be added together with fragrance
// TODO: remake this so all perfumers need to be updated - add 'updated' column
func updatePerfumers(db *config.Database, numRequests int) {
	db.Logger.Info("updating perfumers")
	// fmt.Println("-Updating perfumers-")

	err := fragrances.UpdatePerfumers(db, numRequests)
	if err != nil {
		return
	}
}
