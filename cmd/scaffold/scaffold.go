package scaffold

import (
	"degit/internal/cli"
	"degit/internal/scaffold/config"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var cwd string
var definitions cli.MapVar
var flagInit bool
var dryRun bool
var verbose bool
var showHelp bool

func printHelp() {
	fmt.Println("Usage: degit scaffold [options] <alias>")
	fmt.Println("Options:")
	fmt.Println("   -D <name>=<value>   Define a variable")
	fmt.Println("   --init              Initialize scaffold configuration")
	fmt.Println("   -n                  Dry run - show what would be done without making changes")
	fmt.Println("   -v                  Verbose output")
	fmt.Println("   -h                  Show help")
}

func initFlag() error {
	flag.Var(&definitions, "D", "Define a variable")
	flag.BoolVar(&flagInit, "init", false, "Initialize scaffold configuration")
	flag.BoolVar(&dryRun, "n", false, "Dry run - show what would be done without making changes")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&showHelp, "h", false, "Show help information")
	flag.Usage = printHelp

	var err error
	cwd, err = os.Getwd()
	if err != nil {
		return err
	}
	return nil
}

func Execute(globalPrintHelp func(), die func(error)) {
	if err := initFlag(); err != nil {
		die(err)
	}

	os.Args = os.Args[1:]
	flag.Parse()
	args := flag.Args()

	if showHelp || (!flagInit && len(args) == 0) {
		printHelp()
		return
	}

	if flagInit {
		if err := initialize(); err != nil {
			die(err)
		}
		return
	}

	config, err := config.ReadFile(getConfigFilePath())
	if err != nil {
		die(err)
	}

	if len(args) == 1 && args[0] == "." {
		scaffolder := NewProjectScaffolder(config)
		scaffolder.DryRun = dryRun
		scaffolder.Verbose = verbose
		if err := scaffolder.Deploy(definitions); err != nil {
			die(err)
		}
	}
}

func initialize() error {
	path := getConfigFilePath()
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("configuration file already exists")
	}

	if !os.IsNotExist(err) {
		return err
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(getInitConfigContent())
	if err != nil {
		return err
	}

	return nil
}
