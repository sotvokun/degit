package template

import (
	"fmt"
	"strings"
)

func getOptionExtensions() []string {
	ext, ok := options["extensions"]
	if !ok {
		return []string{}
	}
	return strings.Split(ext, ",")
}

func getOptionDelimiter() ([]string, error) {
	delimiter, ok := options["delimiter"]
	if !ok {
		return []string{"{{", "}}"}, nil
	}
	commaCount := strings.Count(delimiter, ",")
	if commaCount != 1 {
		return nil, fmt.Errorf("invalid delimiter format")
	}
	if delimiter[0] == ',' || delimiter[len(delimiter)-1] == ',' {
		return nil, fmt.Errorf("invalid delimiter format")
	}
	return strings.Split(delimiter, ","), nil
}

func getOptionNonstrict() bool {
	nonstrict, ok := options["nonstrict"]
	if !ok {
		return false
	}
	return strings.ToLower(nonstrict) == "true"
}

func getOptionRemovesource() bool {
	removesource, ok := options["removesource"]
	if !ok {
		return false
	}
	return strings.ToLower(removesource) == "true"
}
