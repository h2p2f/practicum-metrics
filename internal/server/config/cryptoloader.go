package config

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

// cryptoLoader - функция загрузки крипто ключа
//
// cryptoLoader - function of loading crypto key
func (config *ServerConfig) cryptoLoader() {
	config.Logger.Debug("try to load private RSA key")
	if config.KeyFile != "" {
		data, err := os.ReadFile(config.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		block, _ := pem.Decode(data)
		if block == nil {
			log.Fatal("failed to parse PEM block containing the key")
		}
		config.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		config.Logger.Debug("private RSA key loaded successfully")
	} else {
		config.Logger.Debug("private RSA key not loaded")
		config.PrivateKey = nil
	}
}
