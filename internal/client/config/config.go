package config

import "time"

type Config struct {
	Address        string
	ReportInterval time.Duration
	PoolInterval   time.Duration
}

func NewConfig() *Config {
	return &Config{
		Address:        "localhost:8080",
		ReportInterval: 10,
		PoolInterval:   2,
	}
}

func (c *Config) SetConfig(address string, reportInterval time.Duration, poolInterval time.Duration) *Config {
	c.Address = address
	c.ReportInterval = reportInterval
	c.PoolInterval = poolInterval
	return c
}
