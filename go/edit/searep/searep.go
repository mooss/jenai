// Package searep implement the search and replacement block protocol.
package searep

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/mooss/bagend/go/obj/oset"
)

var (
	// SearchMarker marks the start of a search and replace block.
	SearchMarker = "<<<<<<< SEARCH"
	// SeparatorMarker marks the boundary between search and replace.
	SeparatorMarker = "======="
	// ReplaceMarker marks the end of a search and replace block.
	ReplaceMarker = ">>>>>>> REPLACE"
)

// Block is a search and replace block.
type Block struct {
	// Filepath as it appears in the model's response.
	Filepath string
	// Content to search for.
	Search string
	// Replacement content.
	Replace string
}

// ApplySearchAndReplace applies the given search/replace blocks to the editable files and returns
// the resulting in-memory buffers.
//
//revive:disable:cognitive-complexity
func ApplySearchAndReplace(subs []Block, writeFiles []string) (UpdatedFiles, error) {
	updated := make(UpdatedFiles)
	allowed := oset.New(writeFiles...)

	for _, sub := range subs {
		if !allowed.Contains(sub.Filepath) {
			return nil, fmt.Errorf("file %q is not in writable files", sub.Filepath)
		}

		content, ok := updated[sub.Filepath]
		if !ok {
			data, err := os.ReadFile(sub.Filepath)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("when reading editable file %q: %w", sub.Filepath, err)
			}
			content = string(data)
		}

		if sub.Search == "" && content != "" {
			msg := "search content for %q must not be empty for non-empty files"
			return nil, fmt.Errorf(msg, sub.Filepath)
		}

		if count := strings.Count(content, sub.Search); count != 1 {
			msg := "search content for %q occurs %d times; expected exactly once"
			return nil, fmt.Errorf(msg, sub.Filepath, count)
		}

		updated[sub.Filepath] = strings.Replace(content, sub.Search, sub.Replace, 1)
	}

	return updated, nil
}
