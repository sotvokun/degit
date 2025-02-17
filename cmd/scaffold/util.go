package scaffold

import (
	"bufio"
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed init.yaml
var initConfigContent string

func getConfigFilePath() string {
	return filepath.Join(cwd, ".degit", "degit.yaml")
}

func getInitConfigContent() string {
	return initConfigContent
}

func scanString(scanner ...*bufio.Scanner) (string, error) {
	_scanner := bufio.NewScanner(os.Stdin)
	if len(scanner) > 0 {
		_scanner = scanner[0]
	}

	_scanner.Scan()
	if err := _scanner.Err(); err != nil {
		return "", err
	}
	return _scanner.Text(), nil
}
