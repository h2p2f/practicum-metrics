package config

import "time"

type Config struct {
	Address        string
	Key            string
	ReportInterval time.Duration
	PoolInterval   time.Duration
	Batch          bool
	RateLimit      int
}

func NewConfig() *Config {
	return &Config{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PoolInterval:   2 * time.Second,
		Key:            "",
		Batch:          true,
		RateLimit:      6,
	}
}

func (c *Config) SetConfig(address string,
	key string,
	reportInterval time.Duration,
	poolInterval time.Duration,
	batch bool,
	rateLimit int) *Config {
	c.Address = address
	c.ReportInterval = reportInterval
	c.PoolInterval = poolInterval
	c.Key = key
	c.Batch = batch
	c.RateLimit = rateLimit
	return c
}
