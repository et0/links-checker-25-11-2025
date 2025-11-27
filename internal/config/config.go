package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP    HTTP    `yaml:"http_server"`
	Storage Storage `yaml:"storage"`
	App     App     `yaml:"app"`
}

type HTTP struct {
	Port string `yaml:"port"`
}

type Storage struct {
	Filepath string `yaml:"filepath"`
}

type App struct {
	WorkerCount int `yaml:"worker_count"`
}

func Loading() (*Config, error) {
	configPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists || configPath == "" {
		return nil, fmt.Errorf("Wrong config file")
	}

	cfg := &Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	return cfg, nil
}
