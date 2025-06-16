package config

import (
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
	BaseUrl string
	ApiKey  string
}

func NewConfig() *Config {
	file, err := os.ReadFile("config.toml")
	if err != nil {
		return &Config{}
	}

	var config Config
	if err := toml.Unmarshal(file, &config); err != nil {
		return &Config{}
	}
	return &config
}
