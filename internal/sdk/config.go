package sdk

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type StringOrArray []string

func (s *StringOrArray) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		trimmed := strings.TrimSpace(value.Value)
		if trimmed != "" {
			*s = []string{trimmed}
		}
		return nil
	}

	if value.Kind == yaml.SequenceNode {
		var items []string
		if err := value.Decode(&items); err != nil {
			return err
		}
		*s = items
		return nil
	}

	return fmt.Errorf("expected string or array, got %v", value.Kind)
}

type Config struct {
	Mode        string        `yaml:"mode"`
	Dir         string        `yaml:"dir"`
	URL         string        `yaml:"url"`
	FailOn      string        `yaml:"fail_on"`
	Ignore      []string      `yaml:"ignore"`
	Only        []string      `yaml:"only"`
	StartCmd    StringOrArray `yaml:"start_cmd"`
	WaitFor     StringOrArray `yaml:"wait_for"`
	WaitTimeout time.Duration `yaml:"wait_timeout"`
	Format      string        `yaml:"format"`
	Output      string        `yaml:"output"`
	MaxPages    int           `yaml:"max_pages"`
}

func DefaultConfig() Config {
	return Config{
		FailOn:      "high",
		Format:      "text",
		MaxPages:    1000,
		WaitTimeout: 30 * time.Second,
	}
}

func LoadConfig(configPath string) (Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.FailOn == "" {
		config.FailOn = "high"
	}
	if config.Format == "" {
		config.Format = "text"
	}
	if config.MaxPages <= 0 {
		config.MaxPages = 1000
	}
	if config.WaitTimeout <= 0 {
		config.WaitTimeout = 30 * time.Second
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Mode == "" {
		return fmt.Errorf("--mode is required (static or full)")
	}

	if c.Mode != "static" && c.Mode != "full" {
		return fmt.Errorf("--mode must be 'static' or 'full', got '%s'", c.Mode)
	}

	if c.Mode == "static" && c.Dir == "" {
		return fmt.Errorf("--dir is required in static mode")
	}

	if c.Mode == "full" && c.URL == "" {
		return fmt.Errorf("--url is required in full mode")
	}

	validSeverities := map[string]bool{
		"critical": true, "high": true, "medium": true, "low": true, "info": true,
	}
	if !validSeverities[c.FailOn] {
		return fmt.Errorf("--fail-on must be one of: critical, high, medium, low, info")
	}

	if c.Format != "text" && c.Format != "json" {
		return fmt.Errorf("--format must be 'text' or 'json'")
	}

	return nil
}

func (c *Config) MergeFlag(name string, value string) {
	if value == "" {
		return
	}

	switch name {
	case "mode":
		c.Mode = value
	case "dir":
		c.Dir = value
	case "url":
		c.URL = value
	case "fail-on":
		c.FailOn = value
	case "ignore":
		c.Ignore = splitAndTrim(value)
	case "only":
		c.Only = splitAndTrim(value)
	case "start-cmd":
		c.StartCmd = splitAndTrim(value)
	case "wait-for":
		c.WaitFor = splitAndTrim(value)
	case "format":
		c.Format = value
	case "output":
		c.Output = value
	}
}

func splitAndTrim(commaSeparated string) []string {
	parts := strings.Split(commaSeparated, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
