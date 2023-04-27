package config

import "time"

// Config for future usage in handlers
type Configurer interface {
	NewConfig() *serverConfig
	SetServerAddress(address string)
	SetStoreInterval(interval int)
	SetPathToStoreFile(path string)
	SetRestore(restore bool)
	GetConfig() *serverConfig
	SetConfig(string, time.Duration, string, bool) *serverConfig
}

type serverConfig struct {
	// address to listen on
	ServerAddress   string
	StoreInterval   time.Duration
	PathToStoreFile string
	Restore         bool
}

// NewConfig is a function that returns a new config
func NewConfig() *serverConfig {
	return &serverConfig{
		ServerAddress:   "localhost:8080",
		StoreInterval:   300,
		PathToStoreFile: "/tmp/devops-metrics-db.json",
		Restore:         true,
	}
}
func (c *serverConfig) SetServerAddress(address string) {
	c.ServerAddress = address
}
func (c *serverConfig) SetStoreInterval(interval int) {
	c.StoreInterval = time.Duration(interval)
}

func (c *serverConfig) SetPathToStoreFile(path string) {
	c.PathToStoreFile = path
}

func (c *serverConfig) SetRestore(restore bool) {
	c.Restore = restore
}

func (c *serverConfig) GetConfig() *serverConfig {
	return c
}

func (c *serverConfig) SetConfig(address string, interval time.Duration, path string, restore bool) *serverConfig {
	c.ServerAddress = address
	c.StoreInterval = interval
	c.PathToStoreFile = path
	c.Restore = restore
	return c
}
