package scaffold

import (
	"degit/cmd/cmdargs"
	"degit/internal/scaffold/config"
	"degit/internal/scaffold/worker"
	"flag"
	"fmt"
	"os"
)

var definitions cmdargs.MapVar
var showHelp bool
var dryRun bool
var outputPath string
var initConfig bool
var cwd string

func printHelp() {
	fmt.Println("Usage: degit scaffold [options] {<alias> | <filepath>}")
	fmt.Println("Options:")
	fmt.Println("   -D <name>=<value>   Define a variable")
	fmt.Println("   -o <outputpath>     Output path of the result")
	fmt.Println("   --init              Initialize scaffold configuration")
	fmt.Println("   -n                  Dry run - show what would be done without making changes")
	fmt.Println("   -h                  Display help information")
}

func initFlag() error {
	_cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	flag.Var(&definitions, "D", "Define a variable")
	flag.BoolVar(&showHelp, "h", false, "Display help information")
	flag.BoolVar(&dryRun, "n", false, "Dry run - show what would be done without making changes")
	flag.StringVar(&outputPath, "o", "", "Output path of the result")
	flag.BoolVar(&initConfig, "init", false, "Initialize scaffold configuration")
	flag.StringVar(&cwd, "cwd", _cwd, "Current working directory")
	flag.Usage = printHelp
	return nil
}

func Execute(globalHelpFunc func(), die func(error)) {
	if err := initFlag(); err != nil {
		die(err)
	}

	os.Args = os.Args[1:]

	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	args := flag.Args()
	if !initConfig && len(args) == 0 {
		printHelp()
		return
	}

	switch {
	case initConfig:
		err := executeWithInit()
		if err != nil {
			die(err)
		}
	case args[0] == ".":
		err := executeWithDot()
		if err != nil {
			die(err)
		}
	}
}

func executeWithInit() error {
	configPath := config.ConfigFilePath(cwd)
	if config.Exists(configPath) {
		return fmt.Errorf("configuration file already exists")
	}

	err := config.Create(configPath)
	if err != nil {
		return err
	}
	return nil
}

func executeWithDot() error {
	configPath := config.ConfigFilePath(cwd)
	if !config.Exists(configPath) {
		return fmt.Errorf("configuration file does not exist")
	}

	cfg, err := config.LoadFile(configPath)
	if err != nil {
		return err
	}

	return worker.ScaffoldProject(cfg)
}
