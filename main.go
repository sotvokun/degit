package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var branch *string

func printHelp() {
	println("Usage: git clone [-b <branch>] <repo> [<dir>]")
}

func init() {
	flag.Usage = printHelp
	branch = flag.String("b", "", "branch to clone")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}

	src := resolveRepositoryUrl(args[0])
	dest := "."
	if len(args) > 1 {
		dest = args[1]
	}

	ref, err := getRepositoryMatchedRef(src, branch)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheExists, err := checkCacheExists(src, ref)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !cacheExists {
		err := cacheRepository(src, ref)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cacheDirFilename, err := getCacheFileName(src, ref)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cachePath := filepath.Join(cacheDir, cacheDirFilename)

	err = extractCache(cachePath, dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func getCacheFileName(src string, ref HashRef) (string, error) {
	name, err := getRepositoryName(src)
	if err != nil {
		return "", err
	}
	filename := name + "-" + ref.Ref.Short() + "-" + ref.Hash + ".tar.gz"
	return strings.Replace(filename, "/", "--", -1), nil
}

func checkCacheExists(src string, ref HashRef) (bool, error) {
	filename, err := getCacheFileName(src, ref)
	if err != nil {
		return false, err
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filepath.Join(cacheDir, filename))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func cacheRepository(src string, ref HashRef) error {
	filename, err := getCacheFileName(src, ref)
	if err != nil {
		return err
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	filepath := filepath.Join(cacheDir, filename)
	_, fs, err := cloneRepository(src, ref)
	if err != nil {
		return err
	}

	archiveWriter, err := createArchive(filepath)
	if err != nil {
		return err
	}

	err = walkRepositoryFilesystem(fs, "/", func(path string, dir string) error {
		file, err := fs.Open(path)
		if err != nil {
			return err
		}
		fileinfo, err := fs.Stat(path)
		if err != nil {
			return err
		}
		archiveWriter.Add(file, &fileinfo, path)
		file.Close()
		return nil
	})
	if err != nil {
		archiveWriter.Close()
		return err
	}

	err = archiveWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func extractCache(path string, dest string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = uncompressArchive(file, dest)
	if err != nil {
		return err
	}
	return nil
}
