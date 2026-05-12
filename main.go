package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/config"
)

func main() {
	// Setup
	db, logCloser, err := config.Setup()
	if err != nil {
		fmt.Errorf("Error reading config: %v", err)
	}
	defer logCloser()

	fmt.Printf("Application running...\n")
	db.Logger.Info("application started",
		"env", db.BuildEnv,
	)

	ScraperService(db, 30)
}
