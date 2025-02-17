package scaffold

import (
	"degit/internal/command"
	"degit/internal/common/filesystem"
	"degit/internal/scaffold/config"
	"degit/internal/template/inspect"
	"degit/internal/template/pathmap"
	"degit/internal/template/render"
	"fmt"
	"slices"
	"strings"
)

type ProjectScaffold struct {
	config  *config.ProjectConfig
	DryRun  bool
	Verbose bool
}

func NewProjectScaffolder(config *config.Config) *ProjectScaffold {
	return &ProjectScaffold{config: &config.Project}
}

func (s *ProjectScaffold) Deploy(definitions map[string]string) error {
	if s.config == nil {
		return nil
	}
	if len(s.config.Glob) == 0 {
		return fmt.Errorf("no glob patterns defined")
	}

	ctx := command.NewContext()
	ctx.Verbose = s.Verbose
	ctx.DryRun = s.DryRun

	globCmd := filesystem.NewGlobCommand(s.config.Glob)
	if err := globCmd.Execute(ctx); err != nil {
		return err
	}

	files, typeOk := ctx.Get(filesystem.ContextKeyGlobResult).(filesystem.ContextKeyGlobResultType)
	if files == nil || !typeOk {
		return fmt.Errorf("failed to collect files from glob")
	}

	vars, err := s.collectVariables(files)
	if err != nil {
		return err
	}
	if definitions == nil {
		definitions = make(map[string]string)
	}
	err = s.inputVariables(vars, &definitions)
	if err != nil {
		return err
	}

	renderCmd := render.NewRenderCommand(files, definitions)
	if err := s.setupRenderCommand(renderCmd); err != nil {
		return err
	}
	if err := renderCmd.Execute(ctx); err != nil {
		return err
	}

	pathmapCmd := pathmap.New(files)
	if err := s.setupPathMapCommand(pathmapCmd); err != nil {
		return err
	}
	if err := pathmapCmd.Execute(ctx); err != nil {
		return err
	}

	return s.operateResult(ctx)
}

func (s *ProjectScaffold) inputVariables(vars []string, result *map[string]string) error {
	for _, v := range vars {
		_, alreadyDefined := (*result)[v]
		if alreadyDefined {
			continue
		}

		def, exists := s.config.Variables[v]

		fmt.Printf("\033[97m%s\033[0m", v)
		if exists {
			if def.Name != "" {
				fmt.Printf(" \033[33m<%s>\033[0m", def.Name)
			}
			if def.Description != "" {
				fmt.Printf(" \033[90m%s\033[0m", def.Description)
			}
			if def.Default != "" {
				fmt.Printf(" (default: %s)", def.Default)
			}
		}
		fmt.Printf(": ")

		input, err := scanString()
		if err != nil {
			return err
		}

		if strings.TrimSpace(input) == "" && def.Default == "" {
			return fmt.Errorf("no input provided for %s", v)
		}
		if strings.TrimSpace(input) == "" && def.Default != "" {
			input = def.Default
		}
		if len(def.Enum) > 0 && !slices.Contains(def.Enum, input) {
			return fmt.Errorf("invalid input for %s, options are: %v", v, strings.Join(def.Enum, ", "))
		}
		(*result)[v] = input
	}
	return nil
}

func (s *ProjectScaffold) collectVariables(files []string) ([]string, error) {
	ctx := command.NewContext()
	ctx.Verbose = s.Verbose
	ctx.DryRun = s.DryRun

	inspectCmd := inspect.New(files, nil)
	if err := inspectCmd.Execute(ctx); err != nil {
		return nil, err
	}

	vars, typeOk := ctx.Get(inspect.ContextKeyExecuteResult).(inspect.ContextKeyExecuteResultType)
	if vars == nil || !typeOk {
		return nil, fmt.Errorf("failed to collect variables from files")
	}

	return vars, nil
}

func (s *ProjectScaffold) setupRenderCommand(renderCmd *render.RenderCommand) error {
	delimiters, err := s.getOptionDelimiter()
	if err != nil {
		return err
	}
	renderCmd.SetDelimiters(delimiters[0], delimiters[1])

	if s.getOptionNonstrict() {
		renderCmd.SetMissingKeyPolicy(render.MissingKeyPolicyDefault)
	}

	return nil
}

func (s *ProjectScaffold) setupPathMapCommand(pathMapCommand *pathmap.PathMapCommand) error {
	pathMapCommand.ExtensionsToRemove = s.getOptionExtensions()
	return nil
}

func (s *ProjectScaffold) operateResult(ctx *command.Context) error {
	renderResult := ctx.Get(render.ContextKeyExecuteResult).(render.ContextKeyExecuteResultType)
	pathmapResult := ctx.Get(pathmap.ContextKeyExecuteResult).(pathmap.ContextKeyExecuteResultType)

	fsCmds := make([]command.Command, 0)
	for src, content := range renderResult {
		dest, ok := pathmapResult[src]
		if !ok {
			dest = src
		}

		writeCmd := filesystem.NewWriteCommand(dest, []byte(content))
		fsCmds = append(fsCmds, writeCmd)

		if s.getOptionRemovesource() && src != dest {
			removeCmd := filesystem.NewRemoveCommand(src)
			fsCmds = append(fsCmds, removeCmd)
		}
	}

	for _, cmd := range fsCmds {
		if err := cmd.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *ProjectScaffold) getOptionExtensions() []string {
	if s.config == nil {
		return make([]string, 0)
	}
	options, exists := s.config.Options["extensions"]
	if !exists {
		return make([]string, 0)
	}
	return strings.Split(options, ",")
}

func (s *ProjectScaffold) getOptionDelimiter() ([]string, error) {
	if s.config == nil {
		return []string{"{{", "}}"}, nil
	}
	options, exists := s.config.Options["delimiter"]
	if !exists {
		return []string{"{{", "}}"}, nil
	}
	commaCount := strings.Count(options, ",")
	if commaCount != 1 {
		return nil, fmt.Errorf("invalid delimiter format")
	}
	if options[0] == ',' || options[len(options)-1] == ',' {
		return nil, fmt.Errorf("invalid delimiter format")
	}

	return strings.Split(options, ","), nil
}

func (s *ProjectScaffold) getOptionNonstrict() bool {
	if s.config == nil {
		return false
	}
	options, exists := s.config.Options["nonstrict"]
	if !exists {
		return false
	}
	return strings.ToLower(options) == "true"
}

func (s *ProjectScaffold) getOptionRemovesource() bool {
	if s.config == nil {
		return false
	}
	options, exists := s.config.Options["removesource"]
	if !exists {
		return false
	}
	return strings.ToLower(options) == "true"
}
