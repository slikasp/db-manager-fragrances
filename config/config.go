package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/slikasp/dbmanfrags/database"
)

const configFileName = "config.json"

type Config struct {
	RemoteDbURL string `json:"remote_db_url"`
}

type State struct {
	DB        *database.Queries
	CurrentID int32
}

func getConfigFilePath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(currentDir, configFileName)
	return path, nil
}

func Read() (Config, error) {
	// Read the config file in user's HOME directory
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Parse and return the Config struct
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
