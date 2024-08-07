package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type HashRef struct {
	Hash string
	Ref  plumbing.ReferenceName
}

func GetRemoteRefs(url_ string, auth ...Auth) ([]HashRef, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url_},
	})

	listOptions := &git.ListOptions{}

	if len(auth) > 0 {
		authMethod, err := NewAuthMethod(auth[0])
		if err != nil {
			return nil, err
		}
		listOptions.Auth = authMethod
	}

	refs, err := remote.List(listOptions)
	if err != nil {
		return nil, err
	}

	hashRefs := make([]HashRef, len(refs))
	for _, ref := range refs {
		hashRefs = append(hashRefs, HashRef{
			Hash: ref.Hash().String(),
			Ref:  ref.Name(),
		})
	}
	return hashRefs, nil
}

func SearchRef(refs []HashRef, search *string) (*HashRef, error) {
	if search == nil || len(strings.TrimSpace(*search)) == 0 {
		return searchDefaultBranch(refs)
	}
	for _, ref := range refs {
		if ref.Hash == *search {
			return &ref, nil
		}
		tagRef := plumbing.NewTagReferenceName(*search)
		branchRef := plumbing.NewBranchReferenceName(*search)
		if ref.Ref == tagRef || ref.Ref == branchRef {
			return &ref, nil
		}
	}

	return nil, fmt.Errorf("no ref found for %s", *search)
}

func searchDefaultBranch(refs []HashRef) (*HashRef, error) {
	for _, ref := range refs {
		if ref.Ref == plumbing.HEAD {
			return &ref, nil
		}
	}
	return nil, fmt.Errorf("no default branch found")
}
