package inspector

import (
	"slices"
	"text/template"
	"text/template/parse"
)

type Inspector struct {
	variables []string
}

func New() *Inspector {
	return &Inspector{
		variables: make([]string, 0),
	}
}

func (i *Inspector) Variables() []string {
	return i.variables
}

func (i *Inspector) Parse(content string, cache ...bool) ([]string, error) {
	t, err := template.New("template").Parse(content)
	if err != nil {
		return nil, err
	}

	variables := make([]string, 0)
	for _, node := range t.Tree.Root.Nodes {
		i.inspectNode(node, &variables)
	}

	if len(cache) > 0 && cache[0] {
		i.variables = append(i.variables, variables...)
	}

	return variables, nil
}

func (i *Inspector) ParseFiles(files []string, cache ...bool) ([]string, error) {
	t, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	variables := make([]string, 0)
	for _, t := range t.Templates() {
		for _, node := range t.Tree.Root.Nodes {
			i.inspectNode(node, &variables)
		}
	}

	if len(cache) > 0 && cache[0] {
		i.variables = append(i.variables, variables...)
	}
	return variables, nil
}

func (i *Inspector) inspectNode(node parse.Node, variables *[]string) {
	switch n := node.(type) {
	case *parse.FieldNode:
		if !slices.Contains(*variables, n.Ident[0]) {
			*variables = append(*variables, n.Ident[0])
		}
	case *parse.ActionNode:
		i.inspectNode(n.Pipe, variables)
	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			for _, arg := range cmd.Args {
				i.inspectNode(arg, variables)
			}
		}
	case *parse.IfNode:
		i.inspectNode(n.Pipe, variables)
		if n.List != nil {
			for _, node := range n.List.Nodes {
				i.inspectNode(node, variables)
			}
		}
		if n.ElseList != nil {
			for _, node := range n.ElseList.Nodes {
				i.inspectNode(node, variables)
			}
		}
	case *parse.RangeNode:
		i.inspectNode(n.Pipe, variables)
		if n.List != nil {
			for _, node := range n.List.Nodes {
				i.inspectNode(node, variables)
			}
		}
		if n.ElseList != nil {
			for _, node := range n.ElseList.Nodes {
				i.inspectNode(node, variables)
			}
		}
	}
}
