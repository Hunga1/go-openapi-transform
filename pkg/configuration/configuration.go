package configuration

import (
	"encoding/json"
	"os"
)

const CONFIG_LOCATION = "config.json"

type Configuration struct {
	IgnoreFiles []string `json:"ignoreFiles"`
	Whitelist []string `json:"whitelist"`
}

func NewConfig() *Configuration {
	var config Configuration
	configFile, err := os.Open(CONFIG_LOCATION)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	jsonParse := json.NewDecoder(configFile)
	jsonParse.Decode(&config)
	return &config
}
