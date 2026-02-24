package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	cfg1, err := Read()
	if err != nil {
		t.Errorf("Failed to read config: %v", err)
	}

	if cfg1.RemoteDbURL != "test_url" && cfg1.CurrentID != 1 {
		t.Errorf("Bad read output: %s", cfg1.RemoteDbURL)
	}

	cfg1.CurrentID = 10

	err = Write(cfg1)
	if err != nil {
		t.Errorf("Failer to write config: %v", err)
	}

	cfg2, err := Read()
	if err != nil {
		t.Errorf("Failed to read config again: %v", err)
	}

	if cfg2.CurrentID != cfg1.CurrentID {
		t.Errorf("Bad output written: %d:%d", cfg1.CurrentID, cfg2.CurrentID)
	}

	// Reset
	cfg2.CurrentID = 1
	Write(cfg2)
}
