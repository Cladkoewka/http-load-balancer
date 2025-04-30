package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port        int      `yaml:"port"`
	BackendURLs []string `yaml:"backendURLs"`
	HealthCheckInterval int `yaml:"healthCheckInterval"`
	RateLimitDefaults RateLimitDefaults `yaml:"rateLimitDefaults"`
}

type RateLimitDefaults struct {
	BucketRefillRate int `yaml:"bucketRefillRate"`
	BucketCapacity int `yaml:"bucketCapacity"`
}

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
