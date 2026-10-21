package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mooss/bagend/go/fun/eager/lie"
)

func validateWriteFiles(files []string, force bool) error {
	if len(files) == 0 {
		return errors.New("at least one --write file is required")
	}

	validator := writeValidator{seen: map[string]struct{}{}, force: force}
	return lie.Validate(validator.validateWriteFile, files)
}

type writeValidator struct {
	seen  map[string]struct{}
	force bool
}

func (wv *writeValidator) validateWriteFile(file string) error {
	file = strings.TrimSpace(file)
	if file == "" {
		return errors.New("--write contains an empty file path")
	}

	resolved, err := resolveWritePath(file)
	if err != nil {
		return err
	}
	if _, ok := wv.seen[resolved]; ok {
		return fmt.Errorf("file %q was specified more than once with --write", file)
	}
	wv.seen[resolved] = struct{}{}

	if err := validateWritePath(resolved); err != nil {
		return err
	}

	if !wv.force {
		if err := validateGitFile(resolved); err != nil {
			return err
		}
	}
	return nil
}

func resolveWritePath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("make path absolute: %w", err)
	}

	resolved, err := filepath.EvalSymlinks(filepath.Clean(absolute))
	if err != nil {
		return "", fmt.Errorf("resolve symlinks: %w", err)
	}
	return filepath.Clean(resolved), nil
}

func validateWritePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access --write file %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("--write path %q is not a regular file", path)
	}

	handle, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("--write file %q is not writable: %w", path, err)
	}
	if err := handle.Close(); err != nil {
		return fmt.Errorf("close --write file %q: %w", path, err)
	}

	return nil
}

func validateGitFile(file string) error {
	ctx := context.Background()
	repo, err := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("--write file %q must be inside a git repository", file)
	}

	absolute, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("resolve --write file %q: %w", file, err)
	}

	root := strings.TrimSpace(string(repo))
	relative, err := filepath.Rel(root, absolute)
	if err != nil {
		return err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("--write file %q is outside the git repository", file)
	}

	tracked, err := exec.CommandContext(ctx, "git", "-C", root, "ls-files", "--", relative).Output()
	if err != nil {
		return fmt.Errorf("check whether %q is tracked by git: %w", file, err)
	}
	if strings.TrimSpace(string(tracked)) == "" {
		return fmt.Errorf("--write file %q is not tracked by git (use --force to allow it)", file)
	}

	var status bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "-C", root, "status", "--porcelain", "--", relative)
	cmd.Stdout = &status
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("check git status for %q: %w", file, err)
	}
	if hasUnstagedChanges(status.String()) {
		return fmt.Errorf("--write file %q has unstaged changes (use --force to allow it)", file)
	}

	return nil
}

func hasUnstagedChanges(status string) bool {
	for line := range strings.SplitSeq(status, "\n") {
		if line == "" {
			continue
		}
		if len(line) < 2 || line[1] != ' ' {
			return true
		}
	}
	return false
}
