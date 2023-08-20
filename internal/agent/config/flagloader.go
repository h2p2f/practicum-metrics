package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// isSet - функция, проверяющая, установлен ли флаг.
//
// isSet is a function that checks if the flag is set.
func isSet(fs *flag.FlagSet, name string) bool {
	set := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}

// flagLoader - функция загрузки конфигурации из флагов
//
// flagLoader - function of loading configuration from flags
func (config *AgentConfig) flagLoader() {
	config.Logger.Debug("Loading config from flags")
	useJSONConfig := false
	var jsonConfigPath string
	jsonConfFS := flag.NewFlagSet("json", flag.ContinueOnError)
	jsonConfFS.StringVar(&jsonConfigPath, "c", "./config/agent.json", "config file")
	jsonConfFS.StringVar(&jsonConfigPath, "config", "./config/agent.json", "config file")
	err := jsonConfFS.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		log.Println(err)
	}
	if isSet(jsonConfFS, "c") || isSet(jsonConfFS, "config") {
		useJSONConfig = true
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		useJSONConfig = true
		jsonConfigPath = envConfig
	}
	if useJSONConfig {
		jsonFile, err := os.ReadFile(jsonConfigPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonFile, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "Report interval")
	fs.DurationVar(&config.PollInterval, "p", config.PollInterval, "Poll interval")
	fs.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	fs.StringVar(&config.Key, "k", config.Key, "Key")
	fs.StringVar(&config.KeyFile, "crypto-key", config.KeyFile, "RSA key file")
	fs.IntVar(&config.RateLimit, "l", config.RateLimit, "Rate limit")
	// парсим флаги
	// parse flags
	err = fs.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		log.Println(err)
	}
	//если ключ не задан - обнуляем его
	//if the key is not set - zero it
	if !isSet(fs, "k") {
		config.Key = ""
	}
	if !isSet(fs, "crypto-key") {
		config.KeyFile = ""
	}

	config.Logger.Debug("Config loaded from flags")
}
