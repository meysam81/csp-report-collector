package main

import (
	"errors"
	"strings"

	"github.com/meysam81/x/config"
)

type Config struct {
	Port     int    `koanf:"port"`
	LogLevel string `koanf:"log-level"`

	RedisHost        string `koanf:"redis.host"`
	RedisPort        int    `koanf:"redis.port"`
	RedisDB          int    `koanf:"redis.db"`
	RedisPassword    string `koanf:"redis.password"`
	RedisSSLRequired bool   `koanf:"redis.ssl-enabled"`

	RateLimitMaxRPS     int     `koanf:"ratelimit.max"`
	RateLimitRefillRate float32 `koanf:"ratelimit.refill"`
}

func (c *Config) Validate() error {
	errs := []string{}

	if c.RedisHost == "" {
		errs = append(errs, "redis host is empty. provide the value using REDIS_HOST env var")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func NewConfig() (*Config, error) {
	defaults := map[string]interface{}{
		"port":                  8080,
		"log-level":             "info",
		"redis.host":            "localhost",
		"redis.port":            6379,
		"redis.db":              0,
		"ratelimit.max":         20,
		"ratelimit.refill":      2.0,
		"allowed-content-types": []string{"application/csp-report", "application/json", "application/reports+json"},
	}

	c := &Config{}
	_, err := config.NewConfig(config.WithDefaults(defaults), config.WithUnmarshalTo(c))
	if err != nil {
		return nil, err
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}
