package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"go.uber.org/zap"
)

// cryptoLoader - function of loading crypto key
func (config *ServerConfig) cryptoLoader(logger *zap.Logger) (err error) {
	logger.Debug("try to load private RSA key")
	if config.HTTP.KeyFile == "" {
		logger.Debug("private RSA key not loaded")
		config.HTTP.PrivateKey = nil
		return errors.New("private RSA key not loaded because key_file param is empty")
	}
	data, err := os.ReadFile(config.HTTP.KeyFile)
	if err != nil {
		logger.Error("failed to read private RSA key", zap.Error(err))
		return err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		logger.Error("failed to parse PEM block containing the key")
		return err
	}
	config.HTTP.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Error("failed to parse PKCS1 private key", zap.Error(err))
		return err
	}
	logger.Debug("private RSA key loaded successfully")
	return nil
}
