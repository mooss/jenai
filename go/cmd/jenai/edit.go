package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/bagend/go/gitw"
	"github.com/mooss/bagend/go/pth"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/edit"
	"github.com/mooss/jenai/go/edit/searep"
)

func runEdit(qry *config.Query, force bool) error {
	paths, err := qry.Context.WritablePaths()
	if err != nil {
		return err
	}

	prepare := func() error {
		return validateWriteFiles(paths, force)
	}
	mkprompt := func() (config.Prompt, error) {
		return edit.BuildPrompt(qry)
	}
	onResponse := func(response string) error {
		blocks, err := searep.Parse(response)
		if err != nil {
			return err
		}
		buffers, err := searep.ApplyBlocks(blocks, paths)
		if err != nil {
			return err
		}
		return buffers.Write()
	}

	return runWithHooks(qry, mkprompt, prepare, onResponse)
}

// validateWriteFiles checks that the requested files are writable and, unless
// force is set, tracked by Git with no unstaged changes.
func validateWriteFiles(files []string, force bool) error {
	if len(files) == 0 {
		return errors.New("at least one writable file is required")
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
		return errors.New("writable files cannot be empty")
	}

	resolved, err := pth.Resolve(file)
	if err != nil {
		return err
	}
	if _, ok := wv.seen[resolved]; ok {
		return fmt.Errorf("file %q was specified more than once", file)
	}
	wv.seen[resolved] = struct{}{}

	if err := pth.CheckWritable(resolved); err != nil {
		return err
	}

	if !wv.force {
		if err := gitw.FileIsClean(resolved); err != nil {
			return fmt.Errorf("file is not git-clean, use --force to allow it: %w", err)
		}
	}
	return nil
}
