package config

import (
	"go.uber.org/zap"
	"net"
)

func (config *AgentConfig) ipLoader() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		config.Logger.Fatal("Failed to get interface addresses", zap.Error(err))
	}
	for _, addr := range addrs {
		//fmt.Println(addr)
		if ip, ok := addr.(*net.IPNet); ok {
			if ip.IP.IsGlobalUnicast() && ip.IP.To4() != nil {
				config.IPaddr = &ip.IP
				break
			}
		}
	}
}
