package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"go.uber.org/zap"
	"os"
)

// cryptoLoader - функция загрузки крипто ключа
//
// cryptoLoader - function of loading crypto key
func (config *ServerConfig) cryptoLoader(logger *zap.Logger) (err error) {
	logger.Debug("try to load private RSA key")
	if config.Params.KeyFile == "" {
		logger.Debug("private RSA key not loaded")
		config.Params.PrivateKey = nil
		return errors.New("private RSA key not loaded because key_file param is empty")
	}
	data, err := os.ReadFile(config.Params.KeyFile)
	if err != nil {
		logger.Error("failed to read private RSA key", zap.Error(err))
		return err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		logger.Error("failed to parse PEM block containing the key")
		return err
	}
	config.Params.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Error("failed to parse PKCS1 private key", zap.Error(err))
		return err
	}
	logger.Debug("private RSA key loaded successfully")
	return nil
}
