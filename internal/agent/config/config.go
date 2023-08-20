// Package config содержит конфигурацию агента. Инициализация конфига происхоит из файла config/agent.yaml.
// Затем конфигурация перезаписывается флагами командной строки и переменными окружения если они заданы.
//
// Package config contains the agent configuration. The config is initialized from the config/agent.yaml file.
// Then the configuration is overwritten by command line flags and environment variables if they are set.
package config

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// AgentConfig - структура, описывающая конфигурацию агента.
//
// AgentConfig - a structure that describes the agent configuration.
type AgentConfig struct {
	ServerAddress  string        `yaml:"server" json:"address"`
	Key            string        `yaml:"key"`
	KeyFile        string        `yaml:"key_file" json:"crypto_key"`
	LogLevel       string        `yaml:"log_Level"`
	RateLimit      int           `yaml:"rate_limit"`
	ReportInterval time.Duration `yaml:"report" json:"report_interval"`
	PollInterval   time.Duration `yaml:"poll" json:"poll_interval"`
	PublicKey      *rsa.PublicKey
	Logger         *zap.Logger
}

// GetConfig - функция, возвращающая конфигурацию агента.
//
// GetConfig is a function that returns the agent configuration.
func GetConfig() *AgentConfig {
	var config AgentConfig
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
	// инициализируем логгер
	// initialize logger
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
	// в данной секции также обрабатывается пользовательский json файл конфигурацией
	// overwrite config with command line flags
	config.flagLoader()
	// перезаписываем конфиг переменными окружения
	// overwrite config with environment variables
	config.envLoader()
	// загружаем ключи
	// load keys
	config.cryptoLoader()

	// возвращаем конфиг
	// return config
	return &config
}
