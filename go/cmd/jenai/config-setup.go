// This file bootstraps and secures the config dir.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mooss/jenai/go/config"
)

func ensureConfig() error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	return ensureConfigYaml()
}

func ensureConfigYaml() error {
	dir, err := config.Root()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(dir, "config.yaml")

	info, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		return writeConfigFile("config.yaml", []byte{}, 0600)
	} else if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("config file %q is not a regular file", cfgPath)
	}

	if info.Mode().Perm() != 0600 {
		if err := os.Chmod(cfgPath, 0600); err != nil {
			return err
		}

		msg := "warning: permissions of %q changed to 0600 to make it non world-readable to " +
			"protect secrets within\n"
		fmt.Fprintf(os.Stderr, msg, cfgPath)
	}

	return nil
}

// ensureConfigDir tests if the directory already exists and creates it if it doesn't.
func ensureConfigDir() error {
	dir, err := config.Root()
	if err != nil {
		return err
	}

	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("config path %q exists but is not a directory", dir)
		}
		return ensureConfigDirPermissions(dir)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	return nil
}

func ensureConfigDirPermissions(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if info.Mode().Perm() == 0700 {
		return nil
	}

	if err := os.Chmod(dir, 0700); err != nil {
		return err
	}

	msg := "warning: permissions of config directory %q changed to 0700 to make it non " +
		"world-readable to protect secrets within\n"
	fmt.Fprintf(os.Stderr, msg, dir)
	return nil
}

func writeConfigFile(path string, data []byte, perm os.FileMode) error {
	dir, err := config.Root()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, path), data, perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
