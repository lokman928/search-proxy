package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server ProxyConfig
	Brave  BraveConfig
}

type ProxyConfig struct {
	Port int
}

type BraveConfig struct {
	BaseUrl   string            `toml:"base_url"`
	ApiKey    string            `toml:"api_key"`
	RateLimit RateLimiterConfig `toml:"rate_limit"`
}

type RateLimiterConfig struct {
	Enable         bool `toml:"enable"`
	MaxConcurrency int  `toml:"max_concurrency"`
	CooldownTime   int  `toml:"cooldown_time"`
}

func NewConfig() *Config {
	file, err := os.ReadFile("config.toml")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return &Config{}
	}

	var config Config
	if err := toml.Unmarshal(file, &config); err != nil {
		fmt.Println("Error parsing config file:", err)
		return &Config{}
	}
	return &config
}
