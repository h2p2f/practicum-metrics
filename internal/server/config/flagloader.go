package config

import (
	"encoding/json"
	"flag"
	"go.uber.org/zap"
	"log"
	"os"
)

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

// flagLoader - function of loading configuration from flags
func (config *ServerConfig) flagLoader(logger *zap.Logger) {
	logger.Debug("Loading config from flags")
	useJSONConfig := false
	var jsonConfigPath string
	jsonConfFS := flag.NewFlagSet("json", flag.ContinueOnError)
	jsonConfFS.StringVar(&jsonConfigPath, "c", "./config/agent.json", "config file")
	jsonConfFS.StringVar(&jsonConfigPath, "config", "./config/agent.json", "config file")
	err := jsonConfFS.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		logger.Error("failed to parse json config flag", zap.Error(err))
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
			logger.Fatal("failed to read json config file", zap.Error(err))
		}
		err = json.Unmarshal(jsonFile, &config)
		if err != nil {
			logger.Fatal("failed to parse json config file", zap.Error(err))
		}
		logger.Info("json config loaded successfully")
		config.Params.jsonLoaded = true
	}

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&config.Params.Address, "a", config.Params.Address, "Server address")
	fs.StringVar(&config.File.Path, "f", config.File.Path, "File path")
	fs.DurationVar(&config.File.StoreInterval, "i", config.File.StoreInterval, "Store interval")
	fs.BoolVar(&config.File.Restore, "r", config.File.Restore, "Restore")
	fs.StringVar(&config.DB.Dsn, "d", config.DB.Dsn, "Database DSN")
	fs.StringVar(&config.Params.Key, "k", config.Params.Key, "Key")
	fs.StringVar(&config.Params.KeyFile, "crypto-key", config.Params.KeyFile, "RSA key file")
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
		config.Params.Key = ""
	}
	if !isSet(fs, "crypto-key") && !config.Params.jsonLoaded && config.Params.LogLevel != "debug" {
		config.Params.KeyFile = ""
	}
}
