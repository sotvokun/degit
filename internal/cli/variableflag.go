package cli

import (
	"fmt"
	"strings"
)

type VariableFlag map[string]string

func (v *VariableFlag) String() string {
	return ""
}

func (v *VariableFlag) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return fmt.Errorf("invalid variable flag: %s", value)
	}
	if *v == nil {
		*v = make(map[string]string)
	}
	(*v)[parts[0]] = strings.Join(parts[1:], "=")
	return nil
}
