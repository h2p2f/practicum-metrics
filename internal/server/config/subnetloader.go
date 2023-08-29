package config

import (
	"go.uber.org/zap"
	"net"
)

func (config *ServerConfig) subnetLoader(logger *zap.Logger) {
	if config.Params.TrustSubnetString != "" {
		_, trustSubnet, err := net.ParseCIDR(config.Params.TrustSubnetString)
		if err != nil {
			logger.Fatal("Failed to parse trust subnet", zap.Error(err))
		}
		config.Params.TrustSubnet = trustSubnet
	} else {
		logger.Debug("No trust subnet provided")
		config.Params.TrustSubnet = nil
	}
}
