package render

import (
	"degit/internal/command"
	"path/filepath"
	"strings"
	"text/template"
)

const ContextKeyExecuteResult command.ContextKey = "rendering-result"

type ContextKeyExecuteResultType = map[string]string

type RenderCommand struct {
	InputFiles          []string
	Data                map[string]string
	missingKeyPolicy    MissingKeyPolicy
	delimiters          []string
	defaultTemplateName string
}

func NewRenderCommand(inputFiles []string, data map[string]string) *RenderCommand {
	return &RenderCommand{
		InputFiles:          inputFiles,
		Data:                data,
		missingKeyPolicy:    MissingKeyPolicyError,
		delimiters:          []string{"{{", "}}"},
		defaultTemplateName: "___template",
	}
}

func (rc *RenderCommand) SetDelimiters(left string, right string) {
	rc.delimiters = []string{left, right}
}

func (rc *RenderCommand) SetMissingKeyPolicy(policy MissingKeyPolicy) {
	rc.missingKeyPolicy = policy
}

func (rc *RenderCommand) SetDefaultTemplateName(name string) {
	rc.defaultTemplateName = name
}

func (rc *RenderCommand) Execute(ctx *command.Context) error {
	tmpl := template.New(rc.defaultTemplateName)
	tmpl, err := tmpl.ParseFiles(rc.InputFiles...)
	if err != nil {
		return err
	}

	tmpl.Option(rc.missingKeyPolicy.String())
	tmpl.Delims(rc.delimiters[0], rc.delimiters[1])

	ctx.Logf("Rendering with missing key policy: %s", rc.missingKeyPolicy.String())
	ctx.Logf("Rendering with delimiters: %s %s", rc.delimiters[0], rc.delimiters[1])

	result := make(ContextKeyExecuteResultType)
	for _, file := range rc.InputFiles {
		var output strings.Builder
		name := filepath.Base(file)
		if err := tmpl.ExecuteTemplate(&output, name, rc.Data); err != nil {
			return err
		}
		result[file] = output.String()

		ctx.Logf("Rendered %s with short name %s:", file, name)
		if ctx.Verbose {
			ctx.Logf("%s", result[file])
		}
	}
	ctx.Set(ContextKeyExecuteResult, result)
	return nil
}
