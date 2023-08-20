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
func (config *AgentConfig) cryptoLoader() {

	if config.KeyFile != "" {
		config.Logger.Debug("Loading public key")
		data, err := os.ReadFile(config.KeyFile)
		if err != nil {
			config.Logger.Fatal("Failed to read public key", zap.Error(err))
		}
		block, _ := pem.Decode(data)
		if block == nil {
			config.Logger.Fatal("failed to parse PEM block containing the key")
		}
		config.PublicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		config.Logger.Debug("Public key loaded")
	} else {
		config.Logger.Debug("No public key provided")
		config.PublicKey = nil
	}
}
