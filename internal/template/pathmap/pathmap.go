package pathmap

import (
	"degit/internal/command"
	"strings"
)

const ContextKeyExecuteResult command.ContextKey = "pathmap-result"

type ContextKeyExecuteResultType = map[string]string

type PathMapCommand struct {
	InputPaths         []string
	ExtensionsToRemove []string
	predefinedDict     map[string]string
}

func New(paths []string, exts ...[]string) *PathMapCommand {
	_exts := make([]string, 0)
	if len(exts) > 0 && exts[0] != nil && len(exts[0]) > 0 {
		_exts = exts[0]
	}

	return &PathMapCommand{
		InputPaths:         paths,
		ExtensionsToRemove: _exts,
		predefinedDict:     make(map[string]string),
	}
}

func (pmc *PathMapCommand) SetPredefinedDict(dict map[string]string) {
	pmc.predefinedDict = dict
}

func (pmc *PathMapCommand) Execute(ctx *command.Context) error {
	result := make(ContextKeyExecuteResultType)
	for _, path := range pmc.InputPaths {
		predefined, ok := pmc.predefinedDict[path]
		if ok {
			result[path] = predefined
			continue
		}
		result[path] = pmc.process(path)
	}
	ctx.Set(ContextKeyExecuteResult, result)
	return nil
}

func (pmc *PathMapCommand) process(source string) string {
	for _, ext := range pmc.ExtensionsToRemove {
		if strings.HasSuffix(source, ext) {
			fullExt := ext
			if !strings.HasPrefix(ext, ".") {
				fullExt = "." + ext
			}
			source = strings.TrimSuffix(source, fullExt)
		}
	}
	return source
}
