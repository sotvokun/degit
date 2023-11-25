package main

import (
	"fmt"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type HashRef struct {
	Hash string
	Ref  plumbing.ReferenceName
}

type WalkCallback func(path string, dir string) error

func getRepositoryMatchedRef(url_ string, search *string) (HashRef, error) {
	DEFAULT_REF := HashRef{
		Hash: "",
		Ref:  plumbing.HEAD,
	}

	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url_},
	})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return DEFAULT_REF, err
	}

	// find the default branch that HEAD points to
	if search == nil || *search == "" {
		for _, ref := range refs {
			if ref.Name() == plumbing.HEAD {
				*search = ref.Target().Short()
				break
			}
		}
		if search == nil || *search == "" {
			return DEFAULT_REF, fmt.Errorf("no default branch found")
		}
	}

	for _, ref := range refs {
		refHash, refName := ref.Hash().String(), ref.Name()
		if refHash == *search {
			return HashRef{Hash: refHash, Ref: refName}, nil
		}
		tagRef := plumbing.NewTagReferenceName(*search)
		branchRef := plumbing.NewBranchReferenceName(*search)
		if ref.Name() == tagRef || ref.Name() == branchRef {
			return HashRef{Hash: refHash, Ref: refName}, nil
		}
	}

	return DEFAULT_REF, fmt.Errorf("no matched ref %s found", *search)
}

func cloneRepository(url_ string, ref HashRef) (*git.Repository, billy.Filesystem, error) {
	fs := memfs.New()
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url_,
		ReferenceName: ref.Ref,
		Depth:         1,
		SingleBranch:  true,
	})
	if err != nil {
		return nil, nil, err
	}
	return repo, fs, nil
}

func walkRepositoryFilesystem(fs billy.Filesystem, path string, cb WalkCallback) error {
	fis, err := fs.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		itemPath := fi.Name()
		if path != "/" {
			itemPath = fmt.Sprintf("%s/%s", path, fi.Name())
		}
		if fi.IsDir() {
			err := walkRepositoryFilesystem(fs, itemPath, cb)
			if err != nil {
				return err
			}
		} else {
			err := cb(itemPath, path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
