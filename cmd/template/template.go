package template

import (
	"degit/internal/template"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var definitions MapVar
var options MapVar
var globs SliceVar
var dryRun bool
var showHelp bool

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
	flag.Var(&definitions, "D", "Define a variable")
	flag.Var(&options, "s", "Set a option")
	flag.Var(&globs, "g", "Add a glob pattern")
	flag.BoolVar(&dryRun, "n", false, "Dry run - show what would be done without making changes")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.Usage = printHelp
}

func Execute(globalHelpFunc func(), die func(error)) {
	initFlag()

	os.Args = os.Args[1:]

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 || showHelp {
		printHelp()
		os.Exit(1)
	}

	if len(globs) > 0 {
		if err := executeWithGlob(); err != nil {
			die(err)
		}
	} else {
		if err := executeWithPath(args); err != nil {
			die(err)
		}
	}
}

func executeWithGlob() error {
	return nil
}

func executeWithPath(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file paths provided")
	}

	src := ""
	dist := ""

	if !filepath.IsAbs(args[0]) {
		var err error
		src, err = filepath.Abs(args[0])
		if err != nil {
			return err
		}
	} else {
		src = args[0]
	}

	if len(args) > 1 {
		dist = args[1]
	} else {
		dist = src
	}
	if !filepath.IsAbs(dist) {
		var err error
		dist, err = filepath.Abs(dist)
		if err != nil {
			return err
		}
	}

	rawContent, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	content := string(rawContent)

	output, err := template.Render(content, definitions)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println(output)
		return nil
	}

	distDir := filepath.Dir(dist)
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		if err := os.MkdirAll(distDir, 0755); err != nil {
			return err
		}
	}

	outputFile, err := os.Create(dist)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}
