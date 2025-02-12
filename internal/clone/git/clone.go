package git

import (
	"degit/internal/clone/filesystem"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type InMemoryFilesystem = filesystem.InMemoryFilesystem
type ReferenceName = plumbing.ReferenceName

func Clone(url string, ref ReferenceName, auth ...Auth) (*git.Repository, *InMemoryFilesystem, error) {
	cloneOptions := &git.CloneOptions{
		URL:               url,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     ref,
	}

	if len(auth) > 0 {
		gitAuth, err := auth[0].Transport()
		if err != nil {
			return nil, nil, err
		}
		cloneOptions.Auth = gitAuth
	}

	memoryStorage := memory.NewStorage()
	fs := memfs.New()
	repo, err := git.Clone(memoryStorage, fs, cloneOptions)
	if err != nil {
		return nil, nil, err
	}
	return repo, filesystem.NewInMemoryFilesystem(fs), nil
}
