package main

import (
	"degit/internal/degit"
	"degit/internal/git"
	"flag"
	"fmt"
	"os"
	"strings"
)

var branch *string
var sshPrivateKeyPath string
var username string
var password string

func printHelp() {
	fmt.Println("Usage: degit [options] <src>[#<ref>] [<dest>]")
	fmt.Println("Options:")
	fmt.Println("  -i <path>  Path to the SSH private key")
	fmt.Println("  -l <name>  Username for authentication")
	fmt.Println("  -p <pass>  Password or personal access token for authentication or SSH private key passphrase")
}

func init() {
	flag.StringVar(&sshPrivateKeyPath, "i", "", "Path to the SSH private key")
	flag.StringVar(&username, "l", "", "Username for the Git repository")
	flag.StringVar(&password, "p", "", "Password or personal access token for the Git repository")
	flag.Usage = printHelp
}

func main() {
	dest := "."
	ref := ""
	src := ""

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}

	if len(args) > 1 {
		dest = args[1]
	}

	if strings.Contains(args[0], "#") {
		split := strings.Split(args[0], "#")
		ref = split[len(split)-1]
		src = strings.Join(split[:len(split)-1], "#")
	} else {
		src = args[0]
	}

	src = degit.ResolveRemoteUrl(src)
	refs, err := git.GetRemoteRefs(src, git.Auth{
		PrivateKeyPath: sshPrivateKeyPath,
		Username:       username,
		Password:       password,
	})
	if err != nil {
		die(err)
	}

	targetRef, err := git.SearchRef(refs, &ref)
	if err != nil {
		die(err)
		os.Exit(1)
	}

	repositoryCache := degit.NewRepositoryCache(src, *targetRef)
	exists, err := repositoryCache.Exists()
	if err != nil {
		die(err)
	}

	if !exists {
		err := repositoryCache.Cache(degit.RepositoryCacheOptions{
			Force:          true,
			Username:       username,
			Password:       password,
			PrivateKeyPath: sshPrivateKeyPath,
		})
		if err != nil {
			die(err)
		}
	}

	err = repositoryCache.Extract(dest)
	if err != nil {
		die(err)
	}

	os.Exit(0)
}

func die(err error) {
	fmt.Println(err)
	os.Exit(1)
}
