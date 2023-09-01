package config

import (
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

// envLoader - function of loading configuration from environment variables
func (config *AgentConfig) envLoader(logger *zap.Logger) {
	logger.Debug("Loading config from environment variables")
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		config.ServerAddress = envServerAddress
	}

	// if the report interval is set in the environment variable - rewrite
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if isNumeric(envReportInterval) {
			envReportInterval += "s"
		}
		config.ReportInterval, _ = time.ParseDuration(envReportInterval)
	}

	// if the poll interval is set in the environment variable - rewrite
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if isNumeric(envPollInterval) {
			envPollInterval += "s"
		}
		config.PollInterval, _ = time.ParseDuration(envPollInterval)
	}

	// if the key is set in the environment variable - rewrite
	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}

	// if the limit is set in the environment variable - rewrite
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if isNumeric(envRateLimit) {
			config.RateLimit, _ = strconv.Atoi(envRateLimit)
		}
	}

	// if the path to the key is set in the environment variable - rewrite
	if envKryptoKey := os.Getenv("CRYPTO_KEY"); envKryptoKey != "" {
		config.KeyFile = envKryptoKey
	}
	logger.Debug("Config loaded from environment variables")
}
