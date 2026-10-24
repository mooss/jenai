// Package config defines jenai's configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func Root() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "jenai"), nil
}

type GlobalFlags struct {
	OutputFormat string
}

// GlobalFlags is the shared destination of global flags.
var Global GlobalFlags
