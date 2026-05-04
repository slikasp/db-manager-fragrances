package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/config"
)

func main() {
	logFile, err := setupLogging("app.log")
	if err != nil {
		log.Fatalf("failed to setup logging: %v", err)
	}
	defer logFile.Close()

	log.Println("---Application started---")
	fmt.Println("---Application started---")

	// Setup
	frags, err := config.Setup()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	log.Println("--Setup complete--")
	fmt.Println("--Setup complete--")

	// ManualDbUpdate(frags)

	ScraperService(frags, 30)

	log.Println("---Application closed---")
	fmt.Println("---Application closed---")
}
