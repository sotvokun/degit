package executor

import (
	"os"
	"path/filepath"
)

func Glob(patterns []string) ([]string, error) {
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}
	return files, nil
}

func ReadFile(path string) (string, error) {
	rawContent, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(rawContent), nil
}
