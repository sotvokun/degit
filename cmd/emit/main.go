package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/sotvokun/emit/internal/command"
	"github.com/sotvokun/emit/internal/pkg/alflag"
	versionService "github.com/sotvokun/emit/internal/service/version"
)

var (
	help    = alflag.Bool("h, help", false)
	version = alflag.Bool("v, version", false)

	commandsMap = make(map[string]command.Command)
	commands    = []command.Command{
		command.NewVersionCommand(),
		command.NewDegitCommand(),
	}
)

func init() {
	for _, command := range commands {
		commandsMap[command.Name()] = command
	}
}

func main() {
	err := alflag.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing arguments: %v\n", err)
		os.Exit(command.ExitCodeInternalError)
	}

	if *help {
		printUsage()
		os.Exit(command.ExitCodeSuccess)
	}

	if *version {
		versionService := versionService.NewVersionService()
		fmt.Fprintln(os.Stdout, versionService.Version())
		os.Exit(command.ExitCodeSuccess)
	}

	if alflag.NArg() == 0 {
		printUsage()
		os.Exit(command.ExitCodeSuccess)
	}

	commandName := alflag.Arg(0)
	cmd, exists := commandsMap[commandName]
	if !exists {
		fmt.Fprintf(os.Stderr, "emit: '%s' is not a valid command\n", commandName)
		os.Exit(command.ExitCodeArgumentError)
	}

	exitCode, err := cmd.Run(alflag.Args()[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode)
	}
	os.Exit(0)
}

func printUsage() {
	fmt.Fprintf(os.Stdout, "%s\n", strings.TrimSpace(usage()))
}
