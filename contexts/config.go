package contexts

import (
  "encoding/json"
  "os"
  "path/filepath"
)

type SecretInfo struct {
  Path string `json:"path"`
  Key string `json:"key"`
}

type Context struct {
  Name string `json:"name"`
  Secrets map[string]SecretInfo `json:"secrets"`
}

type Config struct {
  Contexts []Context `json:"contexts"`
}

func LoadContexts() (*Config, error) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return nil, err
  }

  // Construct the path to the config file in the .cfctl directory
  configFile := filepath.Join(homeDir, ".cfctl", "contexts.json")
  
  file, err := os.Open(configFile)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  
  var config Config
  // Decode the JSON config file
  decoder := json.NewDecoder(file)
  if err := decoder.Decode(&config); err != nil {
    return nil, err
  }

  return &config, nil
}

