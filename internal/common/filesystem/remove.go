package filesystem

import (
	"degit/internal/command"
	"os"
)

type RemoveCommand struct {
	FilePath string
}

func NewRemoveCommand(filePath string) *RemoveCommand {
	return &RemoveCommand{
		FilePath: filePath,
	}
}

func (r *RemoveCommand) Execute(ctx *command.Context) error {
	ctx.Logf("Removing %s", r.FilePath)
	if ctx.DryRun {
		return nil
	}
	return os.Remove(r.FilePath)
}
