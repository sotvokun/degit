package filesystem

import (
	"degit/internal/command"
	"os"
	"path/filepath"
)

type WriteCommand struct {
	FilePath          string
	Data              []byte
	createNotExistDir bool
	defaultDirPerm    os.FileMode
}

func NewWriteCommand(filePath string, data []byte) *WriteCommand {
	return &WriteCommand{
		FilePath:          filePath,
		Data:              data,
		createNotExistDir: true,
		defaultDirPerm:    os.ModePerm,
	}
}

func (w *WriteCommand) SetCreateNotExistDir(createNotExistDir bool) {
	w.createNotExistDir = createNotExistDir
}

func (w *WriteCommand) SetDefaultDirPerm(defaultDirPerm os.FileMode) {
	w.defaultDirPerm = defaultDirPerm
}

func (w *WriteCommand) Execute(ctx *command.Context) error {
	if w.createNotExistDir {
		dir := filepath.Dir(w.FilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			ctx.Logf("Creating directory %s", dir)
			if !ctx.DryRun {
				if err := os.MkdirAll(dir, w.defaultDirPerm); err != nil {
					return err
				}
			}
		}
	}

	ctx.Logf("Writing to %s", w.FilePath)
	if ctx.DryRun {
		return nil
	}

	file, err := os.Create(w.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(w.Data)
	if err != nil {
		return err
	}

	return nil
}
