package util

import (
	"fmt"
	_url "net/url"
	"path"
	"regexp"
	"strings"
)

const GitHubShortenPattern = `^([a-zA-Z0-9\_\.\-]+)\/([a-zA-Z0-9\_\.\-]+)$`

func ResolveUrl(url string) string {
	matched, _ := regexp.MatchString(GitHubShortenPattern, url)
	if !matched {
		return url
	}
	return fmt.Sprintf("https://github.com/%s", url)
}

func RepositoryName(url string) (string, error) {
	parsed, err := _url.Parse(url)
	if err != nil {
		return "", err
	}
	base := path.Base(parsed.Path)
	return strings.TrimSuffix(base, ".git"), nil
}
