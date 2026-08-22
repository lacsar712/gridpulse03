package config

import "time"

const DefaultTimeout = 8 * time.Second

type Config struct {
	Addr           string
	DBPath         string
	DataDir        string
	AdminToken     string
	RequestTimeout time.Duration
	MaxBodyBytes   int64
	PageSize       int
}

func (c *Config) withDefaults() {
	if c.RequestTimeout <= 0 {
		c.RequestTimeout = DefaultTimeout
	}
	if c.MaxBodyBytes <= 0 {
		c.MaxBodyBytes = 1 << 20
	}
	if c.PageSize <= 0 {
		c.PageSize = 20
	}
	if c.PageSize > 200 {
		c.PageSize = 200
	}
}
