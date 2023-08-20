package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// isSet - функция проверки наличия флага
//
// isSet - function of checking the presence of a flag
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
func (config *ServerConfig) flagLoader() {
	config.Logger.Debug("Loading config from flags")
	useJsonConfig := false
	var jsonConfigPath string
	jsonConfFS := flag.NewFlagSet("json", flag.ContinueOnError)
	jsonConfFS.StringVar(&jsonConfigPath, "c", "./config/agent.json", "config file")
	jsonConfFS.StringVar(&jsonConfigPath, "config", "./config/agent.json", "config file")
	err := jsonConfFS.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		log.Println(err)
	}
	if isSet(jsonConfFS, "c") || isSet(jsonConfFS, "config") {
		useJsonConfig = true
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		useJsonConfig = true
		jsonConfigPath = envConfig
	}
	if useJsonConfig {
		jsonFile, err := os.ReadFile(jsonConfigPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonFile, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&config.Address, "a", config.Address, "Server address")
	fs.StringVar(&config.File.Path, "f", config.File.Path, "File path")
	fs.DurationVar(&config.File.StoreInterval, "i", config.File.StoreInterval, "Store interval")
	fs.BoolVar(&config.File.Restore, "r", config.File.Restore, "Restore")
	fs.StringVar(&config.DB.Dsn, "d", config.DB.Dsn, "Database DSN")
	fs.StringVar(&config.Key, "k", config.Key, "Key")
	fs.StringVar(&config.KeyFile, "crypto-key", config.KeyFile, "RSA key file")
	err = fs.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		log.Println(err)
	}
	if isSet(fs, "f") {
		config.File.UseFile = true
	}
	if isSet(fs, "d") {
		config.DB.UsePG = true
	}
	if !isSet(fs, "k") {
		config.Key = ""
	}
	if !isSet(fs, "crypto-key") {
		config.KeyFile = ""
	}
}
