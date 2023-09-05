package config

import (
	"net"

	"go.uber.org/zap"
)

func (config *ServerConfig) subnetLoader(logger *zap.Logger) {
	if config.HTTP.TrustSubnetString != "" {
		_, trustSubnet, err := net.ParseCIDR(config.HTTP.TrustSubnetString)
		if err != nil {
			logger.Fatal("Failed to parse trust subnet", zap.Error(err))
		}
		config.HTTP.TrustSubnet = trustSubnet
	} else {
		logger.Debug("No trust subnet provided")
		config.HTTP.TrustSubnet = nil
	}
}
