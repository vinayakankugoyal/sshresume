package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Tab represents a single tab configuration.
type Tab struct {
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

// Profile contains user profile information.
type Profile struct {
	Name     string `yaml:"name"`
	Email    string `yaml:"email"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
}

// Config represents the application configuration.
type Config struct {
	Profile Profile `yaml:"profile"`
	Tabs    []Tab   `yaml:"tabs"`
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config.
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Profile.Name == "" {
		return fmt.Errorf("profile.name is required")
	}

	if len(c.Tabs) == 0 {
		return fmt.Errorf("at least one tab is required")
	}

	for i, tab := range c.Tabs {
		if tab.Name == "" {
			return fmt.Errorf("tabs[%d].name is required", i)
		}
		if tab.File == "" {
			return fmt.Errorf("tabs[%d].file is required", i)
		}
	}

	return nil
}
