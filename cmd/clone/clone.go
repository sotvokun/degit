package clone

import (
	"degit/internal/clone/cache"
	"degit/internal/clone/git"
	"degit/internal/clone/util"
	"flag"
	"fmt"
	"os"
	"strings"
)

var sshPrivateKeyPath string
var username string
var password string
var showHelp bool

func printHelp() {
	fmt.Println("Usage: degit [options] <src>[#<ref>] [<dest>]")
	fmt.Println("Options:")
	fmt.Println("  -i <path>        Path to the SSH private key")
	fmt.Println("  -l <name>        Username for authentication")
	fmt.Println("  -p <pass>        Password or personal access token for authentication or SSH private key passphrase")
}

func initFlag() {
	flag.StringVar(&sshPrivateKeyPath, "i", "", "Path to the SSH private key")
	flag.StringVar(&username, "l", "", "Username for the Git repository")
	flag.StringVar(&password, "p", "", "Password or personal access token for the Git repository")
	flag.BoolVar(&showHelp, "h", false, "Show help")

	flag.Usage = printHelp
}

func Execute(globalHelpFunc func(), die func(error)) {
	initFlag()

	flag.Parse()

	args := flag.Args()
	if showHelp || len(args) == 0 {
		printHelp()
		os.Exit(1)
	}

	dest := "."
	if len(args) > 1 {
		dest = args[1]
	}

	src, ref := resolveSource(args[0])
	src = util.ResolveUrl(src)

	auth := git.Auth{
		PrivateKey: sshPrivateKeyPath,
		Username:   username,
		Password:   password,
	}
	refs, err := git.RemoteReferences(src, auth)
	if err != nil {
		die(err)
	}

	targetRef, err := util.SearchRef(refs, ref)
	if err != nil {
		die(err)
	}

	cacheExists, path, err := cache.ExistsInfo(src, targetRef.Name().Short(), targetRef.Hash().String())
	if err != nil {
		die(err)
	}
	if cacheExists {
		err := cache.ExtractArchive(path, dest)
		if err != nil {
			die(err)
		}
	}

	repo, fs, err := git.Clone(src, targetRef.Name(), auth)
	if err != nil {
		die(err)
	}

	path, err = cache.Create(repo, fs)
	if err != nil {
		die(err)
	}

	err = cache.ExtractArchive(path, dest)
	if err != nil {
		die(err)
	}

	os.Exit(0)
}

func resolveSource(src string) (string, string) {
	if !strings.Contains(src, "#") {
		return src, ""
	}

	split := strings.SplitN(src, "#", 2)
	if len(split) < 2 {
		return split[0], ""
	}
	return split[0], split[1]
}
