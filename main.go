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

func printHelp() {
	fmt.Println("Usage: degit <src>[#<ref>] [<dest>]")
}

func init() {
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
	refs, err := git.GetRemoteRefs(src)
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
		err := repositoryCache.Cache(true)
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
