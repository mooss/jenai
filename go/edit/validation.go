package edit

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/bagend/go/gitw"
	"github.com/mooss/bagend/go/pth"
)

// ValidateWritable checks that the requested paths are writable or creatable and, unless force is
// set, that the existing ones are tracked by Git with no unstaged changes.
func ValidateWritable(paths []string, force bool) error {
	if len(paths) == 0 {
		return errors.New("at least one writable file is required")
	}

	repo, err := gitw.New(context.Background())
	if err != nil {
		return err
	}

	return lie.Validate((&writeValidator{force, repo}).isValid, paths)
}

type writeValidator struct {
	force bool
	repo  gitw.Repo
}

func (val *writeValidator) isValid(path string) error {
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

	if !val.force {
		if err := val.repo.FileIsClean(path); err != nil {
			return fmt.Errorf("file is not git-clean, use --force to allow it: %w", err)
		}
	}
	return nil
}
