package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/slikasp/dbmanagerfrags/config"
	"github.com/slikasp/dbmanagerfrags/database"
)

func main() {
	// logFile, err := setupLogging("app.log")
	// if err != nil {
	// 	log.Fatalf("failed to setup logging: %v", err)
	// }
	// defer logFile.Close()

	// Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Load the database
	dbtx, err := sql.Open("postgres", cfg.RemoteDbURL)
	dbQueries := database.New(dbtx)

	// Create database struct to be passed to functions
	db := &config.Database{
		Queries: dbQueries,
		Cfg:     &cfg,
	}

	frag, err := db.Queries.GetFragrance(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(frag)
}
