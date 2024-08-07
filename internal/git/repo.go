package git

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type InMemoryFilesystem struct {
	fs billy.Filesystem
}

type WalkCallback func(path string, dir string) error

func Clone(url_ string, ref HashRef, auth ...Auth) (*git.Repository, *InMemoryFilesystem, error) {
	fs := memfs.New()
	cloneOptions := &gogit.CloneOptions{
		URL:               url_,
		ReferenceName:     ref.Ref,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: gogit.DefaultSubmoduleRecursionDepth,
	}

	if len(auth) > 0 {
		authMethod, err := NewAuthMethod(auth[0])
		if err != nil {
			return nil, nil, err
		}
		cloneOptions.Auth = authMethod
	}

	repo, err := gogit.Clone(memory.NewStorage(), fs, cloneOptions)
	if err != nil {
		return nil, nil, err
	}
	return repo, &InMemoryFilesystem{fs}, nil
}

func (fs InMemoryFilesystem) Walk(path string, cb WalkCallback) error {
	filesInfo, err := fs.fs.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fileInfo := range filesInfo {
		itemPath := fileInfo.Name()
		if path != "/" {
			itemPath = filepath.ToSlash(filepath.Join(path, itemPath))
		}
		if fileInfo.IsDir() {
			if err := fs.Walk(itemPath, cb); err != nil {
				return err
			}
		} else {
			if err := cb(itemPath, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (fs InMemoryFilesystem) Lstat(path string) (os.FileInfo, error) {
	return fs.fs.Lstat(path)
}

func (fs InMemoryFilesystem) Open(path string) (io.ReadCloser, error) {
	return fs.fs.Open(path)
}

func (fs InMemoryFilesystem) Readlink(linkPath string) (string, error) {
	return fs.fs.Readlink(linkPath)
}
