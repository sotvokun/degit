package template

import (
	"degit/internal/template/executor"
	"degit/internal/template/renderer"
	"flag"
	"fmt"
	"os"
	"strings"
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
	if (len(globs) == 0 && len(args) == 0) || showHelp {
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
	files, err := executor.Glob(globs)
	if err != nil {
		return err
	}

	renderer, err := createRenderer()
	if err != nil {
		return err
	}
	renderingResult, err := renderer.RenderFiles(files, definitions)
	if err != nil {
		return err
	}

	result := executor.ProcessResult{}
	for filepath, content := range renderingResult {
		result[filepath] = executor.ProcessResultItem{
			Content: content,
			Output:  "",
		}
	}

	exec := executor.New()
	exec.Extension = getOptionExtensions()

	if dryRun {
		exec.PrintOutput(result)
		return nil
	}

	return exec.WriteOutput(result)
}

func executeWithPath(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file paths provided")
	}

	src := args[0]
	dist := ""
	if len(args) > 1 {
		dist = args[1]
	}

	content, err := executor.ReadFile(src)
	if err != nil {
		return err
	}

	renderer, err := createRenderer()
	if err != nil {
		return err
	}
	content, err = renderer.Render(content, definitions)
	if err != nil {
		return err
	}

	result := executor.ProcessResult{
		src: {
			Content: content,
			Output:  dist,
		},
	}
	exec := executor.New()
	exec.Extension = getOptionExtensions()

	if dryRun {
		exec.PrintOutput(result)
		return nil
	}

	return exec.WriteOutput(result)
}

func createRenderer() (*renderer.Renderer, error) {
	r := renderer.New()
	delimiter, err := getOptionDelimiter()
	if err != nil {
		return nil, err
	}
	r.SetDelimiter(delimiter[0], delimiter[1])
	if getOptionNonstrict() {
		r.SetMissingKeyPolicy(renderer.MissingKeyPolicyDefault)
	}
	return r, nil
}

func getOptionExtensions() []string {
	ext, ok := options["extensions"]
	if !ok {
		return []string{}
	}
	return strings.Split(ext, ",")
}

func getOptionDelimiter() ([]string, error) {
	delimiter, ok := options["delimiter"]
	if !ok {
		return []string{"{{", "}}"}, nil
	}
	commaCount := strings.Count(delimiter, ",")
	if commaCount != 1 {
		return nil, fmt.Errorf("invalid delimiter format")
	}
	if delimiter[0] == ',' || delimiter[len(delimiter)-1] == ',' {
		return nil, fmt.Errorf("invalid delimiter format")
	}
	return strings.Split(delimiter, ","), nil
}

func getOptionNonstrict() bool {
	nonstrict, ok := options["nonstrict"]
	if !ok {
		return false
	}
	return strings.ToLower(nonstrict) == "true"
}
