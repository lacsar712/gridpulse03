package config

import (
	"fmt"
	"strings"

	"example.com/gridpulse/internal/validate"
)

func (c Config) Validate() error {
	c.withDefaults()
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("addr is required")
	}
	if strings.TrimSpace(c.DBPath) == "" {
		return fmt.Errorf("db path is required")
	}
	if strings.TrimSpace(c.DataDir) == "" {
		return fmt.Errorf("data dir is required")
	}
	if err := validate.Token(c.AdminToken); err != nil {
		return err
	}
	return nil
}

func (c Config) Normalized() Config {
	c.withDefaults()
	c.DBPath = strings.TrimSpace(c.DBPath)
	c.DataDir = strings.TrimSpace(c.DataDir)
	c.Addr = strings.TrimSpace(c.Addr)
	return c
}

func (c Config) PublicMeta() map[string]any {
	return map[string]any{
		"addr":      c.Addr,
		"page_size": c.PageSize,
		"timeout":   c.RequestTimeout.String(),
	}
}
