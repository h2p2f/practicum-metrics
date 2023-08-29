package config

import (
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

// envLoader - function of loading configuration from environment variables
func (config *ServerConfig) envLoader(logger *zap.Logger) {

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		config.Params.Address = envAddress
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		config.File.Path = envFilePath
		config.File.UseFile = true
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if isNumeric(envStoreInterval) {
			envStoreInterval += "s"
		}
		config.File.StoreInterval, _ = time.ParseDuration(envStoreInterval)
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		envRestore, err := strconv.ParseBool(envRestore)
		if err != nil {
			logger.Debug("failed to parse RESTORE env variable", zap.Error(err))
		}
		config.File.Restore = envRestore
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		config.DB.Dsn = envDSN
		config.DB.UsePG = true
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Params.Key = envKey
	}

}
