package config

import (
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestCryptoLoader(t *testing.T) {
	logger := zaptest.NewLogger(t)
	var config ServerConfig
	config.yamlLoader("../../../config/server.yaml")

	tests := []struct {
		name   string
		file   string
		result bool
	}{
		{
			name:   "Positive test 1",
			file:   "testdata/private.rsa",
			result: true,
		},
		{
			name:   "Negative test 1",
			file:   "testdata/invalid.rsa",
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.HTTP.KeyFile = tt.file
			err := config.cryptoLoader(logger)
			if tt.result {
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if config.HTTP.PrivateKey == nil {
					t.Errorf("private RSA key should not be nil")
				}
			} else {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
				if config.HTTP.PrivateKey != nil {
					t.Errorf("private RSA key should be nil")
				}
			}
		})
	}
}
