package config

import (
	"encoding/json"
	"flag"
	"os"

	"go.uber.org/zap"
)

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

// flagLoader - function of loading configuration from flags
func (config *AgentConfig) flagLoader(logger *zap.Logger) {
	logger.Debug("Loading config from flags")
	useJSONConfig := false
	var jsonConfigPath string
	jsonConfFS := flag.NewFlagSet("json", flag.ContinueOnError)
	jsonConfFS.StringVar(&jsonConfigPath, "c", "./config/agent.json", "config file")
	jsonConfFS.StringVar(&jsonConfigPath, "config", "./config/agent.json", "config file")
	err := jsonConfFS.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		logger.Error("Failed to parse json config flags", zap.Error(err))
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
			logger.Error("Failed to read json config file", zap.Error(err))
			return
		}
		err = json.Unmarshal(jsonFile, &config)
		if err != nil {
			logger.Error("Failed to parse json config file", zap.Error(err))
			return
		}
		config.jsonLoaded = true
	}

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "Report interval")
	fs.DurationVar(&config.PollInterval, "p", config.PollInterval, "Poll interval")
	fs.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	fs.StringVar(&config.Key, "k", config.Key, "Key")
	fs.StringVar(&config.KeyFile, "crypto-key", config.KeyFile, "RSA key file")
	fs.IntVar(&config.RateLimit, "l", config.RateLimit, "Rate limit")
	fs.BoolVar(&config.UseGRPC, "grpc", config.UseGRPC, "Use gRPC")

	// parse flags
	err = fs.Parse(os.Args[1:]) //nolint:errcheck
	if err != nil {
		logger.Error("Failed to parse flags", zap.Error(err))
		return
	}

	//if the key is not set - zero it
	if !isSet(fs, "k") {
		config.Key = ""
	}
	if !isSet(fs, "crypto-key") && !config.jsonLoaded && config.LogLevel != "debug" {
		config.KeyFile = ""
	}
	if !isSet(fs, "grpc") {
		config.UseGRPC = false
	}

	logger.Debug("Config loaded from flags")
}
