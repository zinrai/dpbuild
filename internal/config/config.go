package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Environment struct {
	Distribution string   `yaml:"distribution"`
	Architecture string   `yaml:"architecture"`
	Mirror       string   `yaml:"mirror"`
	Components   []string `yaml:"components"`
}

type Config struct {
	Environments []Environment `yaml:"environments"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".dpbuild", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if len(c.Environments) == 0 {
		return fmt.Errorf("no environments defined in config")
	}

	for i, env := range c.Environments {
		if env.Distribution == "" {
			return fmt.Errorf("environment %d: distribution is required", i+1)
		}
		if env.Architecture == "" {
			return fmt.Errorf("environment %d: architecture is required", i+1)
		}
		if env.Mirror == "" {
			return fmt.Errorf("environment %d: mirror is required", i+1)
		}
		if len(env.Components) == 0 {
			return fmt.Errorf("environment %d: at least one component is required", i+1)
		}
	}

	return nil
}
