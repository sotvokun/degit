package renderer

import (
	"path/filepath"
	"strings"
	"text/template"
)

const DefaultRendererTemplateName = "template"

type Renderer struct {
	missingKeyPolicy MissingKeyPolicy
	delimiter        []string
}

func New() *Renderer {
	return &Renderer{
		missingKeyPolicy: MissingKeyPolicyError,
		delimiter:        []string{"{{", "}}"},
	}
}

func (r *Renderer) SetDelimiter(left string, right string) {
	r.delimiter = []string{left, right}
}

func (r *Renderer) SetMissingKeyPolicy(policy MissingKeyPolicy) {
	r.missingKeyPolicy = policy
}

func (r *Renderer) Render(content string, data any) (string, error) {
	tmpl := template.New(DefaultRendererTemplateName)
	tmpl.Option(r.missingKeyPolicy.String())
	tmpl.Delims(r.delimiter[0], r.delimiter[1])

	tmpl, err := tmpl.Parse(content)
	if err != nil {
		return "", err
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (r *Renderer) RenderFiles(files []string, data any) (map[string]string, error) {
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	tmpl.Option(r.missingKeyPolicy.String())
	tmpl.Delims(r.delimiter[0], r.delimiter[1])

	result := make(map[string]string)
	for _, file := range files {
		var output strings.Builder
		name := filepath.Base(file)
		if err := tmpl.ExecuteTemplate(&output, name, data); err != nil {
			return nil, err
		}
		result[file] = output.String()
	}

	return result, nil
}
