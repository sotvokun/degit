package worker

import (
	"degit/internal/scaffold/config"
	"errors"
)

func ScaffoldProject(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("failed to load configuration")
	}

	return nil
}
