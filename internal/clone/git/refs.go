package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Reference = plumbing.Reference

func RemoteReferences(url string, auth ...Auth) ([]*Reference, error) {
	remoteName := "origin"
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: remoteName,
		URLs: []string{url},
	})

	listOption := &git.ListOptions{}
	if len(auth) > 0 {
		authMethod, err := auth[0].Transport()
		if err != nil {
			return nil, err
		}
		listOption.Auth = authMethod
	}

	return remote.List(listOption)
}
