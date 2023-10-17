package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	cfg = &Config{}
)

func GetConfig() *Config {
	return cfg
}

func Load(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(content, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
