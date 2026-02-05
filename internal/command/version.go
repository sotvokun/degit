package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/sotvokun/emit/internal/pkg/alflag"
	"github.com/sotvokun/emit/internal/service/version"
)

type VersionCommand struct {
	flagset *alflag.FlagSet

	all  *bool
	help *bool
}

func NewVersionCommand() *VersionCommand {
	flagset := alflag.NewFlagSet("version")
	all := flagset.Bool("a, all", false)
	help := flagset.Bool("h, help", false)

	return &VersionCommand{
		flagset: flagset,

		all:  all,
		help: help,
	}
}

func (v *VersionCommand) Name() string {
	return "version"
}

func (v *VersionCommand) Usage() string {
	return "Usage: emit version [-a | --all]"
}

func (v *VersionCommand) Run(args []string) (int, error) {
	if err := v.flagset.Parse(args); err != nil {
		return ExitCodeInternalError, err
	}

	if *v.help {
		fmt.Fprintln(os.Stdout, strings.TrimSpace(v.Usage()))
		return ExitCodeSuccess, nil
	}

	service := version.NewVersionService()
	if *v.all {
		fmt.Fprintf(os.Stdout, "%s-%s-%s\n", service.Version(), service.Commit(), service.Date().Format("20060102150405"))
		return ExitCodeSuccess, nil
	}

	fmt.Fprintln(os.Stdout, service.Version())
	return ExitCodeSuccess, nil
}
