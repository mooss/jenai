// This file bootstraps and secures the config dir.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/prompts"
)

func ensureConfig() error {
	err, alreadyExists := ensureConfigDir()
	if err != nil {
		return err
	}
	if alreadyExists {
		return ensureConfigYaml()
	}

	dir, err := config.ConfigDir()
	if err != nil {
		return err
	}
	fmt.Println("Initialized config file in", dir)

	err = writeConfigFile("models.yaml", models.EmbeddedBytes, 0644)
	if err != nil {
		return err
	}

	err = writeConfigFile("prompts.yaml", prompts.EmbeddedBytes, 0644)
	if err != nil {
		return err
	}

	return ensureConfigYaml()
}

func ensureConfigYaml() error {
	dir, err := config.ConfigDir()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(dir, "config.yaml")

	ferr, exists := fileExists(cfgPath)
	if ferr != nil {
		return ferr
	}
	if !exists {
		return writeConfigFile("config.yaml", []byte{}, 0600)
	}

	info, err := os.Stat(cfgPath)
	if err != nil {
		return err
	}

	if info.Mode().Perm() != 0600 {
		if err := os.Chmod(cfgPath, 0600); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "warning: permissions of %q changed to 0600 to make it non world-readable to protect secrets within\n", cfgPath)
	}

	return nil
}

// ensureConfigDir tests if the directory already exists and creates it if it doesn't.
// Returns true when the directory already exists.
func ensureConfigDir() (error, bool) {
	dir, err := config.ConfigDir()
	if err != nil {
		return err, false
	}

	err, exists := fileExists(dir)
	if err != nil {
		return err, false
	}

	if exists {
		return ensureConfigDirPermissions(dir), true
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err), false
	}
	return nil, false
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

	fmt.Fprintf(
		os.Stderr,
		"warning: permissions of config directory %q changed to 0700 to make it non world-readable to protect secrets within\n",
		dir,
	)
	return nil
}

func writeConfigFile(path string, data []byte, perm os.FileMode) error {
	dir, err := config.ConfigDir()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, path), data, perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

func fileExists(filename string) (error, bool) {
	_, err := os.Stat(filename)
	if err == nil {
		return nil, true
	}

	if os.IsNotExist(err) {
		err = nil
	}

	return err, false
}
