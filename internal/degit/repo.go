package degit

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

const githubShortenPattern = `^([a-zA-Z0-9\_\.\-]+)\/([a-zA-Z0-9\_\.\-]+)$`

func ResolveRemoteUrl(url_ string) string {
	githubShortenMatched, _ := regexp.MatchString(githubShortenPattern, url_)
	if !githubShortenMatched {
		return url_
	}
	return fmt.Sprintf("https://github.com/%s", url_)
}

func GetRepositoryNameByRemoteUrl(url_ string) (string, error) {
	parsedUrl, err := url.Parse(url_)
	if err != nil {
		return "", err
	}
	base := path.Base(parsedUrl.Path)
	if strings.HasSuffix(base, ".git") {
		base = strings.TrimSuffix(base, ".git")
	}
	return base, nil
}
