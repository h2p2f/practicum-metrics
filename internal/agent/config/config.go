// Package config содержит конфигурацию агента. Инициализация конфига происхоит из файла config/agent.yaml.
// Затем конфигурация перезаписывается флагами командной строки и переменными окружения если они заданы.
//
// Package config contains the agent configuration. The config is initialized from the config/agent.yaml file.
// Then the configuration is overwritten by command line flags and environment variables if they are set.
package config

import (
	"crypto/rsa"
	"fmt"
	"net"
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
	LogLevel       string        `yaml:"log_level"`
	RateLimit      int           `yaml:"rate_limit"`
	RetryCount     int           `yaml:"retry_count"`
	RetryWaitTime  time.Duration `yaml:"retry_wait_time"`
	ReportInterval time.Duration `yaml:"report" json:"report_interval"`
	PollInterval   time.Duration `yaml:"poll" json:"poll_interval"`
	jsonLoaded     bool
	PublicKey      *rsa.PublicKey
	Logger         *zap.Logger
	IPaddr         *net.IP
}

// GetConfig is a function that returns the agent configuration.
func GetConfig() (*AgentConfig, *zap.Logger, error) {
	var config AgentConfig
	config.jsonLoaded = false
	// read the default config from the yaml file
	config.yamlLoader()
	// if the log level is info, warn or error
	// (production run) - remove the crypto key from the default configuration
	// in this case, it can be connected by the launch flag
	// or environment variable
	if config.LogLevel == "info" || config.LogLevel == "warn" || config.LogLevel == "error" {
		config.KeyFile = ""
	}
	// initialize logger
	atom, err := zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom))
	defer logger.Sync() //nolint:errcheck
	fmt.Println(config.KeyFile)
	// overwrite config with command line flags
	config.flagLoader(logger)

	// overwrite config with environment variables
	config.envLoader(logger)

	// load keys
	config.cryptoLoader(logger)

	// put IP address in config
	config.ipLoader()
	fmt.Println(config.LogLevel)

	// return config
	return &config, logger, nil
}
