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
	alreadyExists, err := ensureConfigDir()
	if err != nil {
		return err
	}
	if alreadyExists {
		return ensureConfigYaml()
	}

	dir, err := config.Root()
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
	dir, err := config.Root()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(dir, "config.yaml")

	exists, ferr := fileExists(cfgPath)
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

		msg := "warning: permissions of %q changed to 0600 to make it non world-readable to " +
			"protect secrets within\n"
		fmt.Fprintf(os.Stderr, msg, cfgPath)
	}

	return nil
}

// ensureConfigDir tests if the directory already exists and creates it if it doesn't.
// Returns true when the directory already exists.
func ensureConfigDir() (bool, error) {
	dir, err := config.Root()
	if err != nil {
		return false, err
	}

	exists, err := fileExists(dir)
	if err != nil {
		return false, err
	}

	if exists {
		return true, ensureConfigDirPermissions(dir)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return false, fmt.Errorf("failed to create config dir: %w", err)
	}
	return false, nil
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

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		err = nil
	}

	return false, err
}
