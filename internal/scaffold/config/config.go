package config

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed embed/degit.yaml
var DefaultConfigContent string

func configPath(cwd string) string {
	return filepath.Join(cwd, ".degit", "degit.yaml")
}

func Exists(cwd string) bool {
	configFilePath := configPath(cwd)
	_, err := os.Stat(configFilePath)
	return err == nil
}

func Initialize(cwd string) error {
	if Exists(cwd) {
		return nil
	}

	configFilePath := configPath(cwd)

	err := os.MkdirAll(filepath.Dir(configFilePath), 0755)
	if err != nil {

		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(DefaultConfigContent)

	return nil
}
