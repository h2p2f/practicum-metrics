// Package config содержит в себе логику работы с конфигурацией сервера
//
// package config contains the logic of working with the server configuration
package config

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ServerConfig - структура конфигурации сервера
//
// ServerConfig - server configuration structure
type ServerConfig struct {
	LogLevel   string            `yaml:"log_level"`
	Address    string            `yaml:"host" json:"address"`
	Key        string            `yaml:"key"`
	KeyFile    string            `yaml:"key_file" json:"crypto_key"`
	DB         DatabaseConfig    `yaml:"database"`
	File       FileStorageConfig `yaml:"file_storage"`
	PrivateKey *rsa.PrivateKey
	Logger     *zap.Logger
}

// FileStorageConfig - структура конфигурации файлового хранилища
//
// FileStorageConfig - file storage configuration structure
type FileStorageConfig struct {
	Path          string        `yaml:"path" json:"store_file"`
	StoreInterval time.Duration `yaml:"flush_interval"`
	Restore       bool          `yaml:"restore" json:"restore"`
	UseFile       bool          `yaml:"use_file"`
}

// DatabaseConfig - структура конфигурации базы данных
//
// DatabaseConfig - database configuration structure
type DatabaseConfig struct {
	Dsn   string `yaml:"dsn" json:"database_dsn"`
	UsePG bool   `yaml:"use_pg"`
}

// GetConfig - функция получения конфигурации сервера, обрабатывает yaml файл, флаги и переменные окружения
//
// GetConfig - function of obtaining the server configuration, processes the yaml file, flags and environment variables
func GetConfig() *ServerConfig {

	var config ServerConfig
	// читаем дефлотный конфиг из yaml файла
	// read the default config from the yaml file
	config.yamlLoader()
	// если уровень логирования info, warn или error
	//(запуск на проде) - убираем крипто ключ из дефолтной конфигурации
	// в этом случае подключить его можно флагом запуска
	// или переменной окружения
	// if the log level is info, warn or error
	// (production run) - remove the crypto key from the default configuration
	// in this case, it can be connected by the launch flag
	// or environment variable
	if config.LogLevel == "info" || config.LogLevel == "warn" || config.LogLevel == "error" {
		config.KeyFile = ""
	}
	// конфигурируем логгер
	// configure logger
	atom, err := zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	config.Logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom))
	defer config.Logger.Sync() //nolint:errcheck
	// перезаписываем конфиг флагами командной строки
	// в данной секции также обрабатывается пользовательский json файл с конфигурацией
	// overwrite config with command line flags
	// this section also processes the user json file with configuration
	config.flagLoader()
	// перезаписываем конфиг переменными окружения
	// overwrite config with environment variables
	config.envLoader()
	// загружаем крипто ключ
	// load crypto key
	config.cryptoLoader()

	return &config

}
