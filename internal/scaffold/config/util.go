package config

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed embed/degit.yaml
var DefaultConfigContent string

// ConfigFilePath returns the path to the degit configuration file (degit.yaml)
// by joining the cwd with ".degit/degit.yaml"
func ConfigFilePath(cwd string) string {
	return filepath.Join(cwd, ".degit", "degit.yaml")
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Create(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(DefaultConfigContent)
	return err
}
