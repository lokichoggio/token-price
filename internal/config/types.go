package config

import (
	"token-price/pkg/log"
)

type Http struct {
	Addr string `yaml:"addr"`
}

type Config struct {
	GinMode string `yaml:"ginMode"`
	ApiKey  string `yaml:"apiKey"`

	Http Http       `yaml:"http"`
	Log  log.Config `yaml:"log"`
}
