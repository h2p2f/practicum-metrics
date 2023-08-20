package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// envLoader - функция загрузки конфигурации из переменных окружения
//
// envLoader - function of loading configuration from environment variables
func (config *ServerConfig) envLoader() {

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		config.Address = envAddress
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
			log.Println(err)
		}
		config.File.Restore = envRestore
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		config.DB.Dsn = envDSN
		config.DB.UsePG = true
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}

}
