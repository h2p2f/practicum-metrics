package config

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"go.uber.org/zap"
)

// cryptoLoader - функция загрузки крипто ключа
// cryptoLoader - crypto key loading function
func (config *AgentConfig) cryptoLoader(logger *zap.Logger) {

	if config.KeyFile != "" {
		logger.Debug("Loading public key")
		data, err := os.ReadFile(config.KeyFile)
		if err != nil {
			logger.Fatal("Failed to read public key", zap.Error(err))
		}
		block, _ := pem.Decode(data)
		if block == nil {
			logger.Fatal("failed to parse PEM block containing the key")
		}
		config.PublicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		logger.Debug("Public key loaded")
	} else {
		logger.Debug("No public key provided")
		config.PublicKey = nil
	}
}
