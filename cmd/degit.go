package main

import (
	"degit/cmd/clone"
	"degit/cmd/scaffold"
	"degit/cmd/template"
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("Usage: degit [options] <src>[#<ref>] [<dest>]")
	fmt.Println("Usage: degit template <src> [<dest>]")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "template":
		template.Execute(printHelp, die)
	case "scaffold":
		scaffold.Execute(printHelp, die)
	default:
		clone.Execute(printHelp, die)
	}
}

func die(err error) {
	fmt.Println(err)
	os.Exit(1)
}
