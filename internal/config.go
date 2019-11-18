package internal

import (
	"os"
	"encoding/json"
)

type Config struct {
	ListenAddress string `json:"listenAddress"`
	Shortcuts map[string]string `json:"shortcuts"`
}
func LoadShortcutsFromConfig(configPath string) (Config, error) {
	var configuration Config
	file, err := os.Open(configPath)
	if err != nil {
		return configuration, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&configuration); err != nil {
		return configuration, err
	}

	if configuration.Shortcuts == nil {
		configuration.Shortcuts = make(map[string]string)
	}


	if _, ok := configuration.Shortcuts["*"]; !ok {
		configuration.Shortcuts["*"] = DefaultSearchProvider
	}

	configuration.Shortcuts["help"] = "/"

	if configuration.ListenAddress == "" {
		configuration.ListenAddress = "127.0.0.1"
	}

	return configuration, nil
}

