package main

import (
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

	schedule := "40 8-23 * * *"

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
}

// Go through all IDs that have no cards and update if they are now present
// This part takes quite a long time so I would run this rarely and only report on the progress
// >40k missing cards, takes almost 3 hours now, yikes!
func checkMissingCards(db *config.Database) {
	// TODO: make this run in parallel to everything else
	db.Logger.Info("checking missing cards")
	// fmt.Println("-Checking missing cards-")

	err := cards.CheckMissingCards(db)
	if err != nil {
		db.Logger.Error("CheckMissingCards",
			"err", err,
		)
		return
	}
}

func checkExistingCards(db *config.Database) {
	// TODO: make this run in parallel to everything else
	db.Logger.Info("checking downloaded cards")
	// fmt.Println("-Checking downloaded cards-")

	err := cards.CheckExistingCards(db)
	if err != nil {
		db.Logger.Error("CheckExistingCards",
			"err", err,
		)
		return
	}
}

// Go throug all IDs after the last found card and look for new cards
func lookForNewCards(db *config.Database) {
	// With low number of cards to check you can hit a patch of no cards that will prevent you form fiding new ones
	// 100 - too low
	// 1000 - a bit too much
	// 300 - a safe option, takes less than 2 minutes if none are found
	db.Logger.Info("looking for new cards")
	// fmt.Println("-Looking for new cards-")

	err := cards.FindNewCards(db, 300)
	if err != nil {
		db.Logger.Error("FindNewCards",
			"err", err,
		)
		return
	}
}

// Go through newly found cards, try decoding them and add new fragrance entries
func addMissingFragrances(db *config.Database) {
	db.Logger.Info("adding missing fragrances")
	// fmt.Println("-Adding missing fragrances-")

	err := fragrances.AddMissingFragrances(db)
	if err != nil {
		db.Logger.Error("AddMissingFragrances",
			"err", err,
		)
		return
	}

}

// Go through new fragrances, find relevant data and update them
func updateFragranceDetails(db *config.Database, numRequests int) {
	db.Logger.Info("adding missing details")
	// fmt.Println("-Adding missing details-")

	err := fragrances.UpdateFragrances(db, numRequests)
	if err != nil {
		db.Logger.Error("UpdateFragrances",
			"err", err,
		)
		return
	}

}

// do not use after initial list of perfumers is complete - new perfumer to be added together with fragrance
// TODO: remake this so all perfumers need to be updated - add 'updated' column
func updatePerfumers(db *config.Database, numRequests int) {
	db.Logger.Info("updating perfumers")
	// fmt.Println("-Updating perfumers-")

	err := fragrances.UpdatePerfumers(db, numRequests)
	if err != nil {
		db.Logger.Error("UpdatePerfumers",
			"err", err,
		)
		return
	}
}
