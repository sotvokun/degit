package template

import (
	"strings"
	_template "text/template"
)

func Render(content string, data any) (string, error) {
	tmpl, err := _template.New("template").Parse(content)
	if err != nil {
		return "", err
	}
	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}
	return output.String(), nil
}
