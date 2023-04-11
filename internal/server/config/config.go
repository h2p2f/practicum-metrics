package config

// Config for future usage in handlers
type Config interface {
	GetServerAddress() string
	SetServerAddress(string)
	NewConfig() *ServerConfig
}

type ServerConfig struct {
	// address to listen on
	ServerAddress string
}

// NewConfig is a function that returns a new config
func NewConfig() *ServerConfig {
	return &ServerConfig{
		ServerAddress: "",
	}
}

// GetServerAddress is a function that returns the server address
func (c *ServerConfig) GetServerAddress() string {
	return c.ServerAddress
}

// SetServerAddress is a function that sets the server address
func (c *ServerConfig) SetServerAddress(address string) {
	c.ServerAddress = address
}
