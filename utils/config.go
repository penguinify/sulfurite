package utils

import (
    "os"
    "encoding/json"
)

type Config struct {
    MacrosPath string
    Version string
    MacroInterpreterVersion string
}

func LoadConfig(configPath string) *Config {
    file, _ := os.ReadFile(configPath)
    var config Config

    json.Unmarshal(file, &config)

    return &config
}
