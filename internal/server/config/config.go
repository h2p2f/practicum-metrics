// package config contains the logic of working with the server configuration
package config

import (
	"crypto/rsa"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ServerConfig - server configuration structure
type ServerConfig struct {
	LogLevel string            `yaml:"log_level"`
	HTTP     HTTPServerParams  `yaml:"http_server"`
	GRPC     GRPCServerParams  `yaml:"grpc_server"`
	DB       DatabaseConfig    `yaml:"database"`
	File     FileStorageConfig `yaml:"file_storage"`
}

// ServerParams - server parameters structure
type HTTPServerParams struct {
	Address           string `yaml:"host" json:"address"`
	Key               string `yaml:"key"`
	KeyFile           string `yaml:"key_file" json:"crypto_key"`
	TrustSubnetString string `yaml:"trust_subnet" json:"trusted_subnet"`
	jsonLoaded        bool
	PrivateKey        *rsa.PrivateKey
	TrustSubnet       *net.IPNet
}

// FileStorageConfig - file storage configuration structure
type FileStorageConfig struct {
	Path          string        `yaml:"path" json:"store_file"`
	StoreInterval time.Duration `yaml:"flush_interval"`
	Restore       bool          `yaml:"restore" json:"restore"`
	UseFile       bool          `yaml:"use_file"`
}

// DatabaseConfig - database configuration structure
type DatabaseConfig struct {
	Dsn   string `yaml:"dsn" json:"database_dsn"`
	UsePG bool   `yaml:"use_pg"`
}

type GRPCServerParams struct {
	Address string `yaml:"host" json:"grpc_address"`
}

// GetConfig - function of obtaining the server configuration, processes the yaml file, flags and environment variables
func GetConfig() (*ServerConfig, *zap.Logger, error) {

	var config ServerConfig
	config.HTTP.jsonLoaded = false
	path := "./config/server.yaml"
	// read the default config from the yaml file
	config.yamlLoader(path)

	// if the log level is info, warn or error
	// (production run) - remove the crypto key from the default configuration
	// in this case, it can be connected by the launch flag
	// or environment variable
	if config.LogLevel == "info" || config.LogLevel == "warn" || config.LogLevel == "error" {
		config.HTTP.KeyFile = ""
	}

	// configure logger
	atom, err := zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom))
	defer logger.Sync() //nolint:errcheck

	// overwrite config with command line flags
	// this section also processes the user json file with configuration
	config.flagLoader(logger)

	// overwrite config with environment variables
	config.envLoader(logger)

	// load crypto key
	err = config.cryptoLoader(logger)
	if err != nil {
		logger.Error("failed to load crypto key", zap.Error(err))
		config.HTTP.PrivateKey = nil
	}
	// load trusted subnets
	config.subnetLoader(logger)

	return &config, logger, nil

}
