package cli

import (
	"fmt"
	"strings"
)

type MapVar map[string]string

func (m *MapVar) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MapVar) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format, expected <name>=<value>")
	}
	if *m == nil {
		*m = make(map[string]string)
	}
	(*m)[parts[0]] = strings.Join(parts[1:], "=")
	return nil
}
