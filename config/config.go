package config

import (
  "os"
  "path/filepath"
  "encoding/json"
)

type VaultConfig struct {
  Token string `json:"token"`
  Host string `json:"host"`
  CAPath string `json:"caPath"`
}

func SaveConfig(token, host, caPath string) error {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return err
  }

  configFile := filepath.Join(homeDir, ".cfctl", "vault.json")

  err = os.MkdirAll(filepath.Dir(configFile), 0700)
  if err != nil {
    return err
  }

  config := VaultConfig{
    Token: token,
    Host: host,
    CAPath: caPath,
  }
    
  data, err := json.Marshal(config)
  if err != nil {
    return err
  }

  return os.WriteFile(configFile, data, 0600)
}

func LoadConfig() (*VaultConfig, error) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return nil, err
  }
    
  configFile := filepath.Join(homeDir, ".cfctl", "vault.json")

  data, err := os.ReadFile(configFile)
  if err != nil {
    return nil, err
  }

  var config VaultConfig
  if err = json.Unmarshal(data, &config); err != nil {
    return nil, err
  }

  return &config, nil
}

