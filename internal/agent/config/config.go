// Package config содержит конфигурацию агента. Инициализация конфига происхоит из файла config/agent.yaml.
// Затем конфигурация перезаписывается флагами командной строки и переменными окружения если они заданы.
//
// Package config contains the agent configuration. The config is initialized from the config/agent.yaml file.
// Then the configuration is overwritten by command line flags and environment variables if they are set.
package config

import (
	"crypto/rsa"
	"time"

	"go.uber.org/zap"
)

// AgentConfig - структура, описывающая конфигурацию агента.
//
// AgentConfig - a structure that describes the agent configuration.
type AgentConfig struct {
	ServerAddress  string        `yaml:"server" json:"address"`
	Key            string        `yaml:"key"`
	KeyFile        string        `yaml:"key_file" json:"crypto_key"`
	RateLimit      int           `yaml:"rate_limit"`
	ReportInterval time.Duration `yaml:"report" json:"report_interval"`
	PollInterval   time.Duration `yaml:"poll" json:"poll_interval"`
	PublicKey      *rsa.PublicKey
	Logger         *zap.Logger
}

// GetConfig - функция, возвращающая конфигурацию агента.
//
// GetConfig is a function that returns the agent configuration.
func GetConfig(log *zap.Logger) *AgentConfig {
	var config AgentConfig
	config.Logger = log
	config.yamlLoader()
	config.flagLoader()
	config.envLoader()
	config.cryptoLoader()

	// возвращаем конфиг
	// return config
	return &config
}
