package config

import (
	"os"
	"strconv"
	"time"
)

func (config *AgentConfig) envLoader() {
	config.Logger.Debug("Loading config from environment variables")
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		config.ServerAddress = envServerAddress
	}
	// если интервал отчета задан в переменной окружения - перезаписываем
	// if the report interval is set in the environment variable - rewrite
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if isNumeric(envReportInterval) {
			envReportInterval += "s"
		}
		config.ReportInterval, _ = time.ParseDuration(envReportInterval)
	}
	// если интервал опроса задан в переменной окружения - перезаписываем
	// if the poll interval is set in the environment variable - rewrite
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if isNumeric(envPollInterval) {
			envPollInterval += "s"
		}
		config.PollInterval, _ = time.ParseDuration(envPollInterval)
	}
	// если ключ задан в переменной окружения - перезаписываем
	// if the key is set in the environment variable - rewrite
	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}
	// если лимит задан в переменной окружения - перезаписываем
	// if the limit is set in the environment variable - rewrite
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		if isNumeric(envRateLimit) {
			config.RateLimit, _ = strconv.Atoi(envRateLimit)
		}
	}

	if envKryptoKey := os.Getenv("CRYPTO_KEY"); envKryptoKey != "" {
		config.KeyFile = envKryptoKey
	}
	config.Logger.Debug("Config loaded from environment variables")
}
