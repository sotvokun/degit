package template

import (
	"degit/internal/cli"
	"degit/internal/command"
	"degit/internal/common/filesystem"
	"degit/internal/template/pathmap"
	"degit/internal/template/render"
	"flag"
	"fmt"
	"os"
	"strings"
)

type RenderResult = render.ContextKeyExecuteResultType
type PathMapResult = pathmap.ContextKeyExecuteResultType

var definitions cli.MapVar
var options cli.MapVar
var globMode bool
var dryRun bool
var verbose bool
var showHelp bool
var ctx *command.Context

func printHelp() {
	fmt.Println("Usage: degit template [options] {<filepath> [<resultpath>] | <glob-pattern>+}")
	fmt.Println("Options:")
	fmt.Println("   -D <name>=<value>   Define a variable")
	fmt.Println("   -s <name>=<value>   Set a option")
	fmt.Println("   -g                  Enable glob mode (all arguments will be treated as glob patterns)")
	fmt.Println("   -n                  Dry run - show what would be done without making changes")
	fmt.Println("   -v                  Verbose output")
	fmt.Println("   -h                  Show help")
}

func initFlag() {
	flag.Var(&definitions, "D", "Define a variable")
	flag.Var(&options, "s", "Set a option")
	flag.BoolVar(&globMode, "g", false, "Enable glob mode (all arguments will be treated as glob patterns)")
	flag.BoolVar(&dryRun, "n", false, "Dry run - show what would be done without making changes")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
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

	ctx = command.NewContext()
	ctx.DryRun = dryRun
	ctx.Verbose = verbose

	if globMode {
		if err := executeWithGlob(args); err != nil {
			die(err)
		}
	} else {
		if err := executeWithPath(args); err != nil {
			die(err)
		}
	}
}

func executeWithGlob(patterns []string) error {
	glob := filesystem.NewGlobCommand(patterns)
	if err := glob.Execute(ctx); err != nil {
		return err
	}

	files := ctx.Get(filesystem.ContextKeyGlobResult).(filesystem.ContextKeyGlobResultType)
	if len(files) == 0 {
		return fmt.Errorf("no files found")
	}

	renderCmd := render.NewRenderCommand(files, definitions)
	if err := setupRenderCommand(renderCmd); err != nil {
		return err
	}
	if err := renderCmd.Execute(ctx); err != nil {
		return err
	}

	pathmapCmd := pathmap.New(files)
	if err := setupPathMapCommand(pathmapCmd); err != nil {
		return err
	}
	if err := pathmapCmd.Execute(ctx); err != nil {
		return err
	}

	return operateResult(ctx)
}

func executeWithPath(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file paths provided")
	}

	src := args[0]
	dest := ""
	if len(args) > 1 {
		dest = args[1]
	}

	files := []string{src}

	renderCmd := render.NewRenderCommand(files, definitions)
	if err := setupRenderCommand(renderCmd); err != nil {
		return err
	}
	if err := renderCmd.Execute(ctx); err != nil {
		return err
	}

	pathmapCmd := pathmap.New(files)
	if err := setupPathMapCommand(pathmapCmd); err != nil {
		return err
	}
	if strings.TrimSpace(dest) != "" {
		pathmapCmd.SetPredefinedDict(map[string]string{src: dest})
	}
	if err := pathmapCmd.Execute(ctx); err != nil {
		return err
	}

	return operateResult(ctx)
}

func setupRenderCommand(renderCommand *render.RenderCommand) error {
	delimiters, err := getOptionDelimiter()
	if err != nil {
		return err
	}
	renderCommand.SetDelimiters(delimiters[0], delimiters[1])

	if getOptionNonstrict() {
		renderCommand.SetMissingKeyPolicy(render.MissingKeyPolicyDefault)
	}

	return nil
}

func setupPathMapCommand(pathMapCommand *pathmap.PathMapCommand) error {
	pathMapCommand.ExtensionsToRemove = getOptionExtensions()
	return nil
}

func operateResult(ctx *command.Context) error {
	renderResult := ctx.Get(render.ContextKeyExecuteResult).(RenderResult)
	pathmapResult := ctx.Get(pathmap.ContextKeyExecuteResult).(PathMapResult)

	filesystemCommands := make([]command.Command, 0)
	for src, content := range renderResult {
		dest, ok := pathmapResult[src]
		if !ok {
			dest = src
		}

		writeCmd := filesystem.NewWriteCommand(dest, []byte(content))
		filesystemCommands = append(filesystemCommands, writeCmd)

		if getOptionRemovesource() && src != dest {
			removeCmd := filesystem.NewRemoveCommand(src)
			filesystemCommands = append(filesystemCommands, removeCmd)
		}
	}

	for _, cmd := range filesystemCommands {
		if err := cmd.Execute(ctx); err != nil {
			return err
		}
	}

	return nil
}
