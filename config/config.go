package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GeneralModel string `json:"general_model"`
	CodeModel    string `json:"code_model"`
	// ...add other config fields as needed...
}

func DefaultConfig() Config {
	return Config{
		GeneralModel: "ollama_general",
		CodeModel:    "ollama_code",
	}
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return DefaultConfig(), err
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

func SaveConfig(path string, cfg Config) error {
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}
