package config

import (
	"testing"
)

func TestConfigSetup(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Errorf("Failed to read config: %v", err)
	}

	if db.BuildEnv != "dev" {
		t.Errorf("Wrong env variable, got: %s, expected: dev", db.BuildEnv)
	}
}
