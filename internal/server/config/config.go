// Package config содержит в себе логику работы с конфигурацией сервера
//
// package config contains the logic of working with the server configuration
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

// ServerConfig - структура конфигурации сервера
//
// ServerConfig - server configuration structure
type ServerConfig struct {
	LogLevel string            `yaml:"log_level"`
	Address  string            `yaml:"host"`
	Key      string            `yaml:"key"`
	DB       DatabaseConfig    `yaml:"database"`
	File     FileStorageConfig `yaml:"file_storage"`
}

// FileStorageConfig - структура конфигурации файлового хранилища
//
// FileStorageConfig - file storage configuration structure
type FileStorageConfig struct {
	Path          string        `yaml:"path"`
	StoreInterval time.Duration `yaml:"flush_interval"`
	Restore       bool          `yaml:"restore"`
	UseFile       bool          `yaml:"use_file"`
}

// DatabaseConfig - структура конфигурации базы данных
//
// DatabaseConfig - database configuration structure
type DatabaseConfig struct {
	Dsn   string `yaml:"dsn"`
	UsePG bool   `yaml:"use_pg"`
}

// GetConfig - функция получения конфигурации сервера, обрабатывает yaml файл, флаги и переменные окружения
//
// GetConfig - function of obtaining the server configuration, processes the yaml file, flags and environment variables
func GetConfig() *ServerConfig {

	var config *ServerConfig

	file, err := os.Open("./config/server.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err2 := file.Close(); err2 != nil {
			log.Println(err2)
		}
	}()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	config = &ServerConfig{
		LogLevel: "info",
		Address:  "localhost:8080",
		File: FileStorageConfig{
			Path:          "/tmp/metrics-db.json",
			StoreInterval: 10 * time.Second,
			Restore:       true,
			UseFile:       false,
		},
		DB: DatabaseConfig{
			Dsn:   "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			UsePG: false,
		},
		Key: "secret",
	}

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&config.Address, "a", config.Address, "Server address")
	fs.StringVar(&config.File.Path, "f", config.File.Path, "File path")
	fs.DurationVar(&config.File.StoreInterval, "i", config.File.StoreInterval, "Store interval")
	fs.BoolVar(&config.File.Restore, "r", config.File.Restore, "Restore")
	fs.StringVar(&config.DB.Dsn, "d", config.DB.Dsn, "Database DSN")
	fs.StringVar(&config.Key, "k", config.Key, "Key")
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
	return config

}

// isNumeric - функция проверки строки на наличие в ней только цифр
//
// isNumeric - function of checking a string for the presence of only numbers in it
func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

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
