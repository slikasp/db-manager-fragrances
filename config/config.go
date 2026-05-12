package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/slikasp/dbmanfrags/database"
)

// const configFileName = "config.json"

type Database struct {
	Queries  *database.Queries
	BuildEnv string
	Logger   *slog.Logger
}

func Setup() (*Database, func() error, error) {
	// Read config
	err := godotenv.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("No .env file found")
	}
	build := "dev"
	if os.Getenv("BUILD_ENV") != "" {
		build = os.Getenv("BUILD_ENV")
	}
	db := Database{
		BuildEnv: build,
	}
	postgresURL := os.Getenv("POSTGRES_URL")

	// Load the database
	dbtx, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, nil, err
	}
	db.Queries = database.New(dbtx)

	logger, logCloser, err := initialiseLogging("app.log")
	if err != nil {
		return nil, nil, err
	}

	db.Logger = logger

	return &db, logCloser, nil
}

func initialiseLogging(logFile string) (*slog.Logger, func() error, error) {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}
	closer := func() error {
		return file.Close()
	}

	fileHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	consoleHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	logger := slog.New(slog.NewMultiHandler(
		fileHandler,
		consoleHandler,
	))

	return logger, closer, nil
}
