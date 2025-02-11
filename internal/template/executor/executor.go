package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProcessResultItem struct {
	Content string
	Output  string
}

type ProcessResult = map[string]ProcessResultItem

type Executor struct {
	Extension    []string
	RemoveSource bool
}

func New() *Executor {
	return &Executor{
		Extension:    []string{},
		RemoveSource: false,
	}
}

func (e *Executor) WriteOutput(result ProcessResult) error {
	for source, r := range result {
		if strings.TrimSpace(r.Output) == "" {
			r.Output = e.processOutput(source)
		}

		outputDir := filepath.Dir(r.Output)
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				return err
			}
		}

		outputFile, err := os.Create(r.Output)
		if err != nil {
			return err
		}

		_, err = outputFile.WriteString(r.Content)
		if err != nil {
			outputFile.Close()
			return err
		}

		err = outputFile.Close()
		if err != nil {
			return err
		}

		if e.RemoveSource && source != r.Output {
			if err := os.Remove(source); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Executor) PrintOutput(result ProcessResult) {
	for source, r := range result {
		if strings.TrimSpace(r.Output) == "" {
			r.Output = e.processOutput(source)
		}

		format := "\033[32m%s\033[0m \033[1m<-\033[0m "
		if e.RemoveSource && source != r.Output {
			format += "\033[31m%s [REMOVED]\033[0m\n"
		} else {
			format += "\033[33m%s\033[0m\n"
		}
		fmt.Printf(format, r.Output, source)
		fmt.Println(r.Content)
		if r.Content[len(r.Content)-1] != '\n' {
			fmt.Println()
		}
	}
}

func (e *Executor) processOutput(source string) string {
	if len(e.Extension) == 0 {
		return source
	}
	for _, ext := range e.Extension {
		if strings.HasSuffix(source, ext) {
			fullExt := ext
			if !strings.HasPrefix(ext, ".") {
				fullExt = "." + ext
			}
			return strings.TrimSuffix(source, fullExt)
		}
	}
	return source
}
