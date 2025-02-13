package scaffold

import (
	"degit/cmd/cmdargs"
	"flag"
	"fmt"
	"os"
)

var definitions cmdargs.MapVar
var showHelp bool
var dryRun bool
var outputPath string
var initConfig bool

func printHelp() {
	fmt.Println("Usage: degit scaffold [options] {<alias> | <filepath>}")
	fmt.Println("Options:")
	fmt.Println("   -D <name>=<value>   Define a variable")
	fmt.Println("   -o <outputpath>     Output path of the result")
	fmt.Println("   --init              Initialize scaffold configuration")
	fmt.Println("   -n                  Dry run - show what would be done without making changes")
	fmt.Println("   -h                  Display help information")
}

func initFlag() {
	flag.Var(&definitions, "D", "Define a variable")
	flag.BoolVar(&showHelp, "h", false, "Display help information")
	flag.BoolVar(&dryRun, "n", false, "Dry run - show what would be done without making changes")
	flag.StringVar(&outputPath, "o", "", "Output path of the result")
	flag.BoolVar(&initConfig, "init", false, "Initialize scaffold configuration")
	flag.Usage = printHelp
}

func Execute(globalHelpFunc func(), die func(error)) {
	initFlag()

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

	_, err := os.Getwd()
	if err != nil {
		die(err)
	}
}
