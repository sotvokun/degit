package template

import (
	"path/filepath"
	"strings"
	_template "text/template"
)

func Render(content string, data any) (string, error) {
	tmpl, err := _template.New("template").Parse(content)
	if err != nil {
		return "", err
	}
	tmpl.Option("missingkey=error")
	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}
	return output.String(), nil
}

func GlobRender(globs []string, data any) (map[string]string, error) {
	files := make([]string, 0)
	for _, g := range globs {
		matches, err := filepath.Glob(g)
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}

	tmpl, err := _template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	tmpl.Option("missingkey=error")

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
