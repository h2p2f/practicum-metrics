package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	ServerAddress  string        `yaml:"server"`
	ReportInterval time.Duration `yaml:"report"`
	PollInterval   time.Duration `yaml:"poll"`
	Key            string        `yaml:"key"`
	RateLimit      int           `yaml:"rate_limit"`
}

func GetConfig() *AgentConfig {
	var config *AgentConfig

	file, err := os.Open("./config/agent.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "Report interval")
	fs.DurationVar(&config.PollInterval, "p", config.PollInterval, "Poll interval")
	fs.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	fs.StringVar(&config.Key, "k", config.Key, "Key")
	fs.IntVar(&config.RateLimit, "l", config.RateLimit, "Rate limit")

	fs.Parse(os.Args[1:]) //nolint:errcheck

	if !IsSet(fs, "k") {
		config.Key = ""
	}

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		config.ServerAddress = envServerAddress
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if isNumeric(envReportInterval) {
			envReportInterval += "s"
		}
		config.ReportInterval, _ = time.ParseDuration(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if isNumeric(envPollInterval) {
			envPollInterval += "s"
		}
		config.PollInterval, _ = time.ParseDuration(envPollInterval)
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if isNumeric(envRateLimit) {
			config.RateLimit, _ = strconv.Atoi(envRateLimit)
		}
	}

	return config
}

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func IsSet(fs *flag.FlagSet, name string) bool {
	set := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}
