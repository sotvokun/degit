package cache

import (
	"degit/internal/clone/archive"
	"degit/internal/clone/filesystem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
)

func cacheDir() (string, error) {
	appDir := "degit"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(cacheDir, "Temp", appDir), nil
	}
	return filepath.Join(cacheDir, appDir), nil
}

func Create(repo *git.Repository, fs *filesystem.InMemoryFilesystem) (string, error) {
	cacheDirPath, err := cacheDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cacheDirPath, os.ModePerm); err != nil {
		return "", err
	}

	fileName, err := repoCacheName(*repo)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(cacheDirPath, fileName)
	if _, err := os.Lstat(filePath); err == nil {
		return filePath, nil
	}

	aw, err := archive.NewWriter(filePath)
	if err != nil {
		return "", err
	}
	defer aw.Close()

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
			if err := aw.File(file, fileinfo, path); err != nil {
				return err
			}
			return nil
		}

		if fileinfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			target, err := fs.Readlink(path)
			if err != nil {
				return err
			}
			if err := aw.Symlink(fileinfo, path, target); err != nil {
				return err
			}
			return nil
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func Exists(repo *git.Repository) (bool, string, error) {
	cacheDirPath, err := cacheDir()
	if err != nil {
		return false, "", err
	}

	fileName, err := repoCacheName(*repo)
	if err != nil {
		return false, "", err
	}

	filePath := filepath.Join(cacheDirPath, fileName)
	if _, err := os.Lstat(filePath); err == nil {
		return true, filePath, nil
	}
	return false, "", nil
}

func ExistsInfo(url string, shortenRefName string, hash string) (bool, string, error) {
	cacheDirPath, err := cacheDir()
	if err != nil {
		return false, "", err
	}

	fileName, err := repoCacheNamePlain(url, shortenRefName, hash)
	if err != nil {
		return false, "", err
	}
	filePath := filepath.Join(cacheDirPath, fileName)
	if _, err := os.Lstat(filePath); err == nil {
		return true, filePath, nil
	}
	return false, "", nil
}

func ExtractArchive(path string, dest string) error {
	ar, err := archive.NewReader(path)
	if err != nil {
		return err
	}
	defer ar.Close()

	return ar.Uncompress(dest)
}

func Extract(repo *git.Repository, dest string) error {
	cacheDirPath, err := cacheDir()
	if err != nil {
		return err
	}

	fileName, err := repoCacheName(*repo)
	if err != nil {
		return err
	}

	filePath := filepath.Join(cacheDirPath, fileName)
	if exists, _, err := Exists(repo); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("cache for repository does not exist")
	}

	return ExtractArchive(filePath, dest)
}
