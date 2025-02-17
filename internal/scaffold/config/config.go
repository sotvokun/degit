package config

import (
	_ "embed"
	"os"

	"github.com/goccy/go-yaml"
)

type VariableDefinition struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Default     string   `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type ProjectConfig struct {
	Glob      []string                      `json:"glob,omitempty"`
	Options   map[string]string             `json:"options,omitempty"`
	Variables map[string]VariableDefinition `json:"variables,omitempty"`
}

type Config struct {
	Project ProjectConfig `json:".,omitempty"`
}

func ReadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
