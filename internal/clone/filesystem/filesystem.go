package filesystem

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

type WalkCallback = func(path string, dir string) error

type InMemoryFilesystem struct {
	billy.Filesystem
}

func NewInMemoryFilesystem(fs billy.Filesystem) *InMemoryFilesystem {
	return &InMemoryFilesystem{fs}
}

// Walk traverses the file tree rooted at the given path, calling the callback function
// for each file (not directory) in the tree, including the root directory.
//
// For each file, the callback function is called with two arguments:
//   - path: the complete relative path of the file
//   - dir: the directory path containing the file
//
// If the callback function returns an error, the walk stops and returns that error.
// The files are walked in lexical order, which makes the output deterministic.
//
// Example:
//
//	filesystem.Walk("/", func(path, dir string) error {
//	    fmt.Printf("File: %s in directory: %s\n", path, dir)
//	    return nil
//	})
func (imfs InMemoryFilesystem) Walk(path string, cb WalkCallback) error {
	// Read all files and subdirectories in the current directory
	filesInfo, err := imfs.Filesystem.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fileInfo := range filesInfo {
		// Get the file or directory name
		itemPath := fileInfo.Name()

		// Path handling logic:
		// 1. If current path is root ("/"), keep the filename as is
		// 2. Otherwise, join current path with filename to create full path
		// Examples:
		//   - When path="/" : itemPath="foo.txt"
		//   - When path="/dir" : itemPath="dir/foo.txt"
		if path != "/" {
			itemPath = filepath.ToSlash(filepath.Join(path, itemPath))
		}

		if fileInfo.IsDir() {
			// If it's a directory, recursively walk through its contents
			if err := imfs.Walk(itemPath, cb); err != nil {
				return err
			}
		} else {
			// If it's a file, call the callback function with:
			// arg1: itemPath - complete relative path of the file
			// arg2: path - directory path containing the file
			if err := cb(itemPath, path); err != nil {
				return err
			}
		}
	}
	return nil
}

// Lstat returns a FileInfo describing the named file. If the file is
// a symbolic link, the returned FileInfo describes the symbolic link.
// Lstat makes no attempt to follow the link.
func (imfs InMemoryFilesystem) Lstat(path string) (os.FileInfo, error) {
	return imfs.Filesystem.Lstat(path)
}

// Open opens the named file for reading. If successful, methods on the returned
// file can be used for reading; the associated file descriptor has mode O_RDONLY.
func (imfs InMemoryFilesystem) Open(path string) (io.ReadCloser, error) {
	return imfs.Filesystem.Open(path)
}

// Readlink returns the target path of link
func (imfs InMemoryFilesystem) Readlink(link string) (string, error) {
	return imfs.Filesystem.Readlink(link)
}
