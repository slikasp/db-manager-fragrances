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

func Setup() (*Database, error) {
	// Read config
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("No .env file found")
	}
	db := Database{
		BuildEnv: os.Getenv("BUILD_ENV"),
	}
	postgresURL := os.Getenv("POSTGRES_URL")

	// Load the database
	dbtx, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	db.Queries = database.New(dbtx)

	logger, logCloser, err := initialiseLogging("app.log")
	if err != nil {
		return nil, err
	}
	defer logCloser()

	db.Logger = logger

	return &db, nil
}

func initialiseLogging(logFile string) (*slog.Logger, func() error, error) {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}
	closer := func() error {
		return file.Close()
	}

	debugHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	errorHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})

	logger := slog.New(slog.NewMultiHandler(
		debugHandler,
		errorHandler,
	))

	return logger, closer, nil
}
