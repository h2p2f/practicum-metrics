package config

import "time"

type Config struct {
	Address        string
	Key            string
	ReportInterval time.Duration
	PoolInterval   time.Duration
}

func NewConfig() *Config {
	return &Config{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PoolInterval:   2 * time.Second,
		Key:            "",
	}
}

func (c *Config) SetConfig(address string, key string, reportInterval time.Duration, poolInterval time.Duration) *Config {
	c.Address = address
	c.ReportInterval = reportInterval
	c.PoolInterval = poolInterval
	c.Key = key
	return c
}
