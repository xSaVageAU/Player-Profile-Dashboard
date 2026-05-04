package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NATS   NATSConfig   `yaml:"nats"`
	Server ServerConfig `yaml:"server"`
}

type NATSConfig struct {
	URL            string `yaml:"url"`
	Token          string `yaml:"token"`
	Bucket         string `yaml:"bucket"`
	EconomyBucket  string `yaml:"economy_bucket"`
	EconomySubject string `yaml:"economy_subject"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
