package degit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/sotvokun/emit/internal/service/log"
)

type Walker struct {
	dest   string
	logger log.Logger

	dryMode bool
}

func NewWalker(dest string) *Walker {
	return &Walker{
		dest: dest,
	}
}

func (w *Walker) SetLogger(logger log.Logger) {
	w.logger = logger
}

func (w *Walker) SetDryMode(dryMode bool) {
	w.dryMode = dryMode
}

func (w *Walker) WalkCopy(path string, fs billy.Filesystem) error {
	fi, err := fs.Lstat(path)
	if err != nil {
		return err
	}

	destFullPath := filepath.Join(w.dest, path)
	destFullPathDir := filepath.Dir(destFullPath)
	if err := w.do(func() error { return os.MkdirAll(destFullPathDir, os.ModePerm) }); err != nil {
		return err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		link, err := fs.Readlink(path)
		if err != nil {
			return err
		}

		w.log("create symlink: %s -> %s", destFullPath, link)
		return w.do(func() error {
			return os.Symlink(link, destFullPath)
		})
	}

	if _, err := os.Lstat(destFullPath); err == nil {
		return fmt.Errorf("%s: %w", destFullPath, os.ErrExist)
	}

	w.log("create file: %s", destFullPath)
	return w.do(func() error {
		dstFile, err := os.Create(destFullPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		srcFile, err := fs.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		if _, err = io.Copy(dstFile, srcFile); err != nil {
			return err
		}
		return nil
	})
}

func (w *Walker) log(format string, a ...any) {
	if w.logger == nil {
		return
	}
	w.logger.Printf(format+"\n", a...)
}

func (w *Walker) do(fn func() error) error {
	if w.dryMode {
		return nil
	}
	return fn()
}
