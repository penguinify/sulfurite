package utils

import (
    "os"
    "encoding/json"
)

type ConfigJSON struct {
    MacrosPath string
    Version string
    MacroInterpreterVersion string
}


func LoadConfig(configPath string) *ConfigJSON {
    file, _ := os.ReadFile(configPath)
    var config ConfigJSON

    json.Unmarshal(file, &config)

    return &config
}
