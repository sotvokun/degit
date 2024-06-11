package storage

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetCacheDir(dir ...string) (string, error) {
	appCacheDir := "degit"
	if len(dir) > 0 {
		appCacheDir = dir[0]
	}
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(cacheHome, "Temp", appCacheDir), nil
	}
	return filepath.Join(cacheHome, appCacheDir), nil
}
