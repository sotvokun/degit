package command

const (
	ExitCodeSuccess = iota
	ExitCodeArgumentError
	ExitCodeInternalError
)

type Command interface {
	Name() string
	Usage() string
	Run(args []string) (int, error)
}
