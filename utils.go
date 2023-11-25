package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func resolveRepositoryUrl(url_ string) string {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\_\.\-]+\/[a-zA-Z0-9\_\.\-]+$`, url_)
	if !matched {
		return url_
	}
	return fmt.Sprintf("https://github.com/%s", url_)
}

func getRepositoryName(url_ string) (string, error) {
	parsedUrl, err := url.Parse(url_)
	if err != nil {
		return "", err
	}
	base := path.Base(parsedUrl.Path)
	if len(path.Ext(base)) > 0 {
		parts := strings.Split(base, ".")
		return strings.Join(parts[:len(parts)-1], "."), nil
	} else {
		return base, nil
	}
}

func getCacheDir() (string, error) {
	const DIR = "degit"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(cacheDir, "Temp", DIR), nil
	} else {
		return filepath.Join(cacheDir, DIR), nil
	}
}
