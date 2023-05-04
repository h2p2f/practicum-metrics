package config

import "time"

// Configurer interface for config
type Configurer interface {
	NewConfig() *serverConfig
	GetConfig() *serverConfig
	SetConfig(string, time.Duration, string, bool) *serverConfig
}

// serverConfig is a struct that contains server config
type serverConfig struct {
	// address to listen on
	ServerAddress string
	// interval to store metrics
	StoreInterval time.Duration
	// path to store file
	PathToStoreFile string
	// restore metrics from file
	Restore bool

	Database string
}

// NewConfig is a function that returns a new config
func NewConfig() *serverConfig {
	return &serverConfig{
		ServerAddress:   "localhost:8080",
		StoreInterval:   300,
		PathToStoreFile: "/tmp/devops-metrics-db.json",
		Restore:         true,
		Database:        "host=localhost user=practicum password=yandex dbname=postgres sslmode=disable",
	}
}

// SetServerAddress is a function that sets server address

// GetConfig is a function that returns config
func (c *serverConfig) GetConfig() *serverConfig {
	return c
}

// SetConfig is a function that sets config
func (c *serverConfig) SetConfig(address string, interval time.Duration, path string, restore bool, db string) *serverConfig {
	c.ServerAddress = address
	c.StoreInterval = interval
	c.PathToStoreFile = path
	c.Restore = restore
	c.Database = db
	return c
}
