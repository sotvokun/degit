package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var branch *string

func init() {
	flag.Usage = printHelp
	branch = flag.String("b", string(plumbing.HEAD), "branch to clone")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}

	var src string = resolveRepoUrl(args[0])
	dest, err := getDir(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(args) > 1 {
		dest = args[1]
	}
	dest = normalizePath(dest)
	if !isEmptyOrNoneExistent(dest) {
		fmt.Printf("destination path '%s' already exists and is not an empty directory.\n", dest)
		os.Exit(1)
	}

	_, err = git.PlainClone(dest, false, &git.CloneOptions{
		URL:           src,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName(*branch),
		SingleBranch:  true,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = pruneRepository(dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func printHelp() {
	fmt.Println("Usage: degit [-b <branch>] <src> [<dest>]")
}

func resolveRepoUrl(url_ string) string {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\_\.\-]+\/[a-zA-Z0-9\_\.\-]+$`, url_)
	if !matched {
		return url_
	}
	return fmt.Sprintf("https://github.com/%s", url_)
}

func getDir(url_ string) (string, error) {
	u, err := url.Parse(url_)
	if err != nil {
		return "", err
	}
	base := path.Base(u.Path)
	if len(path.Ext(base)) == 0 {
		return path.Base(u.Path), nil
	} else {
		return strings.Split(base, ".")[0], nil
	}
}

func normalizePath(path_ string) string {
	if strings.HasPrefix(path_, "~/") {
		dirname, _ := os.UserHomeDir()
		return path.Join(dirname, path_[2:])
	}

	if strings.HasPrefix(path_, "./") {
		dirname, _ := os.Getwd()
		return path.Join(dirname, path_[2:])
	}

	return path_
}

func checkRepository(path_ string) (bool, error) {
	_, err := os.Stat(path.Join(path_, ".git"))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func pruneRepository(path_ string) error {
	exists, err := checkRepository(path_)
	if !exists && err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("not a git repository: %s", path_)
	}
	return os.RemoveAll(path.Join(path_, ".git"))
}

// Write a function that check the given path, return true if it is not exist or empty, return false if it is not empty.
func isEmptyOrNoneExistent(path_ string) bool {
	fstat, err := os.Stat(path_)
	if err != nil {
		return os.IsNotExist(err)
	}

	if !fstat.IsDir() {
		return false
	}

	f, err := os.Open(path_)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdir(1)
	return err == io.EOF
}
