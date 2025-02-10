package template

import (
	"flag"
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("Usage: degit template [options] <filepath> [<resultpath>]")
	fmt.Println("Options:")
	fmt.Println("   -D <name>=<value>   Define a variable")
	fmt.Println("   -s <name>=<value>   Set a option")
	fmt.Println("   -g <glob>           Add a glob pattern (<filepath> and <resultpath> will be ignored)")
	fmt.Println("   -n                  Dry run - show what would be done without making changes")
	fmt.Println("   -h                  Show help")
}

func initFlag() {
	flag.Usage = printHelp
}

func Execute(globalHelpFunc func(), die func(error)) {
	initFlag()

	os.Args = os.Args[2:]

	flag.Parse()
}
