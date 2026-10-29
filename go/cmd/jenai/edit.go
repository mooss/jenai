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
	resolved, err := qry.Context.Resolve()
	if err != nil {
		return err
	}
	paths := resolved.Editable()

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

// validateWriteFiles checks that the requested paths are writable or creatable and, unless force is
// set, that the existing ones are tracked by Git with no unstaged changes.
func validateWriteFiles(paths []string, force bool) error {
	if len(paths) == 0 {
		return errors.New("at least one writable file is required")
	}

	return lie.Validate((&writeValidator{force}).isValid, paths)
}

type writeValidator struct {
	force bool
}

func (wv *writeValidator) isValid(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("writable files cannot be empty")
	}

	if !pth.Exists(path) {
		_, err := pth.MayCreate(path)
		return err
	}

	if err := pth.CheckWritable(path); err != nil {
		return err
	}

	if !wv.force {
		if err := gitw.FileIsClean(path); err != nil {
			return fmt.Errorf("file is not git-clean, use --force to allow it: %w", err)
		}
	}
	return nil
}
