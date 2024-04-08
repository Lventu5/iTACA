package tests

import (
	"categorizer/config"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cfg, err := config.ParseConfig("../config/config.json")
	if err != nil {
		t.Error("Error parsing config file")
	}

	t.Log(cfg)

	if cfg.Retriever.Type != "Tulip" {
		t.Error("Error parsing Retriever.Type")
	}
	if cfg.Retriever.Host != "localhost" {
		t.Error("Error parsing Retriever.Address")
	}
	if cfg.Retriever.Port != 3000 {
		t.Error("Error parsing Retriever.Port")
	}
}
