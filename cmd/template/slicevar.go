package template

import "fmt"

type SliceVar []string

func (s *SliceVar) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *SliceVar) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0)
	}
	*s = append(*s, value)
	return nil
}
