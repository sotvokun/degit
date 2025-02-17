package filesystem

import (
	"degit/internal/command"
	"path/filepath"
)

const ContextKeyGlobResult command.ContextKey = "glob-result"

type ContextKeyGlobResultType = []string

type GlobCommand struct {
	Patterns []string
}

func NewGlobCommand(patterns []string) *GlobCommand {
	return &GlobCommand{Patterns: patterns}
}

func (g *GlobCommand) Execute(ctx *command.Context) error {
	files := make(ContextKeyGlobResultType, 0)
	for _, pattern := range g.Patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		files = append(files, matches...)

		ctx.Logf("Globbed %d files from pattern %s", len(matches), pattern)
		if ctx.Verbose {
			for _, match := range matches {
				ctx.Logf("  %s", match)
			}
		}
	}
	ctx.Set(ContextKeyGlobResult, files)
	return nil
}
