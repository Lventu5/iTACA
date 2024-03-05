package config

import (
	"encoding/json"
	"io"
	"os"
)

type RetrieverConfig struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

type AnalyserConfig struct {
	Type       string `json:"type"`
	Address    string `json:"address"`
	Port       uint16 `json:"port"`
	Collection string `json:"collection"`
	ApiKey     string `json:"api_key"`
}

type Config struct {
	Retriever RetrieverConfig `json:"retriever"`
	Analyser  AnalyserConfig  `json:"analyser"`
	Log       bool            `json:"log"`
}

// ParseConfig : parses a json file into a Config struct
func ParseConfig(path string) (Config, error) {
	var config Config
	fp, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer fp.Close()

	bytes, err := io.ReadAll(fp)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
