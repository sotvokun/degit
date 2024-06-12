package degit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"degit/internal/archive"
	"degit/internal/git"
	"degit/internal/storage"
)

type RepositoryCache struct {
	remoteUrl string
	ref       git.HashRef

	localFileName *string
}

func NewRepositoryCache(remoteUrl string, ref git.HashRef) *RepositoryCache {
	return &RepositoryCache{
		remoteUrl: remoteUrl,
		ref:       ref,
	}
}

func (c *RepositoryCache) GetLocalFileName() (string, error) {
	if c.localFileName != nil {
		return *c.localFileName, nil
	}
	name, err := GetRepositoryNameByRemoteUrl(c.remoteUrl)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf(
		"%s-%s-%s.tar.gz",
		name,
		c.ref.Ref.Short(),
		c.ref.Hash,
	)
	fileName = strings.Replace(fileName, "/", "--", -1)
	c.localFileName = &fileName
	return fileName, nil
}

func (c *RepositoryCache) Exists() (bool, error) {
	localFileName, err := c.GetLocalFileName()
	if err != nil {
		return false, err
	}
	cacheDir, err := storage.GetCacheDir()
	if err != nil {
		return false, err
	}

	_, err = os.Lstat(filepath.Join(cacheDir, localFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *RepositoryCache) Extract(dest string) error {
	exists, err := c.Exists()
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cache not found")
	}

	cacheDir, err := storage.GetCacheDir()
	if err != nil {
		return err
	}
	localFileName, err := c.GetLocalFileName()
	if err != nil {
		return err
	}

	cachedFilePath := filepath.Join(cacheDir, localFileName)
	file, err := os.Open(cachedFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := archive.Uncompress(file, dest); err != nil {
		return err
	}
	return nil
}

func (c *RepositoryCache) Cache(force ...bool) error {
	doForce := false
	if len(force) > 0 {
		doForce = force[0]
	}

	exists, err := c.Exists()
	if err != nil {
		return err
	}
	if exists && !doForce {
		return nil
	}

	_, fs, err := git.Clone(c.remoteUrl, c.ref)
	if err != nil {
		return err
	}

	cacheDir, err := storage.GetCacheDir()
	if err != nil {
		return err
	}
	localFileName, err := c.GetLocalFileName()
	if err != nil {
		return err
	}
	cachedFilePath := filepath.Join(cacheDir, localFileName)

	a, err := archive.New(cachedFilePath)
	if err != nil {
		return err
	}
	defer a.Close()

	err = fs.Walk("/", func(path string, dir string) error {
		fileinfo, err := fs.Lstat(path)
		if err != nil {
			return err
		}

		if fileinfo.Mode().IsRegular() {
			file, err := fs.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			if err := a.Add(file, fileinfo, path); err != nil {
				return err
			} else {
				return nil
			}
		}

		if fileinfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			target, err := fs.Readlink(path)
			if err != nil {
				return err
			}
			if err := a.Symlink(fileinfo, path, target); err != nil {
				return err
			} else {
				return nil
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
