package inspect

import (
	"degit/internal/command"
	"slices"
	"text/template"
	"text/template/parse"
)

const ContextKeyExecuteResult command.ContextKey = "inspecting-result"

type ContextKeyExecuteResultType = []string

type InspectCommand struct {
	Files    []string
	Contents []string
}

func New(files []string, contents []string) *InspectCommand {
	if files == nil {
		files = make([]string, 0)
	}
	if contents == nil {
		contents = make([]string, 0)
	}
	return &InspectCommand{
		Files:    files,
		Contents: contents,
	}
}

func (i *InspectCommand) Execute(ctx *command.Context) error {
	variables := make([]string, 0)
	if err := i.inspectFiles(ctx, &variables); err != nil {
		return err
	}
	if err := i.inspectContents(ctx, &variables); err != nil {
		return err
	}
	ctx.Set(ContextKeyExecuteResult, variables)
	return nil
}

func (i *InspectCommand) inspectFiles(ctx *command.Context, result *[]string) error {
	if len(i.Files) == 0 {
		return nil
	}
	t, err := template.ParseFiles(i.Files...)
	if err != nil {
		return err
	}

	for _, t := range t.Templates() {
		vars := make([]string, 0)
		for _, node := range t.Tree.Root.Nodes {
			i.traverseNode(node, &vars, *result)
		}
		*result = append(*result, vars...)
		ctx.Logf("Inspecting template %s found %d variables", t.Name(), len(vars))
		if ctx.Verbose {
			for _, v := range vars {
				ctx.Logf("  %s", v)
			}
		}
	}
	return nil
}

func (i *InspectCommand) inspectContents(ctx *command.Context, result *[]string) error {
	for _, content := range i.Contents {
		t, err := template.New("template").Parse(content)
		if err != nil {
			return err
		}

		vars := make([]string, 0)
		for _, node := range t.Tree.Root.Nodes {
			i.traverseNode(node, &vars, *result)
		}
		*result = append(*result, vars...)
		ctx.Logf("Inspecting template string found %d variables", len(vars))
		if ctx.Verbose {
			for _, v := range vars {
				ctx.Logf("  %s", v)
			}
		}
	}
	return nil
}

func (i *InspectCommand) traverseNode(node parse.Node, variables *[]string, globalVars []string) {
	switch n := node.(type) {
	case *parse.FieldNode:
		if !slices.Contains(globalVars, n.Ident[0]) {
			*variables = append(*variables, n.Ident[0])
		}
	case *parse.ActionNode:
		i.traverseNode(n.Pipe, variables, globalVars)
	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			for _, arg := range cmd.Args {
				i.traverseNode(arg, variables, globalVars)
			}
		}
	case *parse.IfNode:
		i.traverseNode(n.Pipe, variables, globalVars)
		if n.List != nil {
			for _, node := range n.List.Nodes {
				i.traverseNode(node, variables, globalVars)
			}
		}
		if n.ElseList != nil {
			for _, node := range n.ElseList.Nodes {
				i.traverseNode(node, variables, globalVars)
			}
		}
	case *parse.RangeNode:
		i.traverseNode(n.Pipe, variables, globalVars)
		if n.List != nil {
			for _, node := range n.List.Nodes {
				i.traverseNode(node, variables, globalVars)
			}
		}
		if n.ElseList != nil {
			for _, node := range n.ElseList.Nodes {
				i.traverseNode(node, variables, globalVars)
			}
		}
	}
}
