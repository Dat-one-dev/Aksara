package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Commands map[string]string `json:"commands"`
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/aksara/config.json"
	}
	return filepath.Join(home, ".config", "aksara", "config.json")
}

func Load() (Config, error) {
	p := Path()
	data, err := os.ReadFile(p)
	if err != nil {
		def := Config{Commands: map[string]string{
			"oc": "opencode",
		}}
		if mkErr := Save(def); mkErr != nil {
			return def, mkErr
		}
		return def, nil
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{Commands: map[string]string{}}, err
	}
	if c.Commands == nil {
		c.Commands = map[string]string{}
	}
	return c, nil
}

func Save(c Config) error {
	p := Path()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}
