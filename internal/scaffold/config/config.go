package config

import (
	"errors"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Project   ProjectConfig             `yaml:".,omitempty"`
	Scaffolds map[string]ScaffoldConfig `yaml:"scaffolds,omitempty"`
}

type VariableDefinition struct {
	FriendlyName string `yaml:"friendlyName,omitempty"`
	Description  string `yaml:"description,omitempty"`
	Example      string `yaml:"example,omitempty"`
}

type ProjectConfig struct {
	Glob        []string                      `yaml:"glob,omitempty"`
	Options     map[string]string             `yaml:"options,omitempty"`
	Variables   map[string]VariableDefinition `yaml:"variables,omitempty"`
	Files       []string                      `yaml:"files,omitempty"`
	Directories []string                      `yaml:"directories,omitempty"`
}

type ScaffoldConfig struct {
	Variables map[string]VariableDefinition `yaml:"variables,omitempty"`
	Mappings  map[string]string             `yaml:"mappings,omitempty"`
}

func LoadFile(path string) (*Config, error) {
	if !Exists(path) {
		return nil, errors.New("file does not exist")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
