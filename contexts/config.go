package contexts

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SecretInfo struct {
	Path string `json:"path"`
	Key  string `json:"key"`
}

type Context struct {
	Name    string                `json:"name"`
	Secrets map[string]SecretInfo `json:"secrets"`
}

type Config struct {
	Contexts []Context `json:"contexts"`
}

func LoadContexts(path string) (*Config, error) {
	if path == "" {
		// If no path provided, use the default location
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(homeDir, ".cfctl", "contexts.json")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
