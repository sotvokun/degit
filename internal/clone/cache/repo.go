package cache

import (
	"degit/internal/clone/util"
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
)

func repoCacheNamePlain(url string, refName string, hash string) (string, error) {
	name, err := util.RepositoryName(url)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf(
		"%s-%s-%s.tar.gz",
		name,
		refName,
		hash,
	)
	fileName = strings.ReplaceAll(fileName, "/", "___")
	return fileName, nil
}

func repoCacheName(repo git.Repository) (string, error) {
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}
	if remote == nil {
		return "", errors.New("only remote repository support to be cached")
	}
	if remote.Config() == nil || remote.Config().URLs == nil || len(remote.Config().URLs) == 0 {
		return "", nil
	}
	remoteUrl := remote.Config().URLs[0]

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	refName := ref.Name().Short()
	hash := ref.Hash().String()

	return repoCacheNamePlain(remoteUrl, refName, hash)
}
