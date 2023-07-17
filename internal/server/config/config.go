package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
)

type ServerConfig struct {
	LogLevel string            `yaml:"log_level"`
	Address  string            `yaml:"host"`
	File     FileStorageConfig `yaml:"file_storage"`
	DB       DatabaseConfig    `yaml:"database"`
	Key      string            `yaml:"key"`
}

type FileStorageConfig struct {
	Path          string        `yaml:"path"`
	StoreInterval time.Duration `yaml:"flush_interval"`
	Restore       bool          `yaml:"restore"`
	UseFile       bool          `yaml:"use_file"`
}

type DatabaseConfig struct {
	Dsn   string `yaml:"dsn"`
	UsePG bool   `yaml:"use_pg"`
}

func GetConfig() *ServerConfig {

	//var config *ServerConfig

	//file, err := os.Open("./config/server.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//decoder := yaml.NewDecoder(file)
	//err = decoder.Decode(&config)
	//if err != nil {
	//	log.Fatal(err)
	//}

	config := &ServerConfig{
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
	fs.Parse(os.Args[1:]) //nolint:errcheck

	if IsSet(fs, "f") {
		config.File.UseFile = true
	}
	if IsSet(fs, "d") {
		config.DB.UsePG = true
	}
	if !IsSet(fs, "k") {
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
