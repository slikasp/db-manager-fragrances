package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanfrags/config"
)

func main() {
	log.Println("---Application started---")
	fmt.Println("---Application started---")

	// Setup
	db, err := config.Setup()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	log.Println("--Setup complete--")
	fmt.Println("--Setup complete--")

	// ManualDbUpdate(frags)

	ScraperService(db, 30)

	log.Println("---Application closed---")
	fmt.Println("---Application closed---")
}
