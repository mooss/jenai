// Package searep implement the search and replacement block protocol.
package searep

import (
	"fmt"
	"os"
	"strings"
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

//revive:disable:cognitive-complexity
func ApplyBlocks(subs []Block, writeFiles []string) (Buffers, error) {
	buffers := make(Buffers)
	allowed := make(map[string]struct{}, len(writeFiles))
	for _, file := range writeFiles {
		allowed[file] = struct{}{}
	}

	for _, sub := range subs {
		_, ok := allowed[sub.Filepath]
		if !ok {
			return nil, fmt.Errorf("file %q is not in writable files", sub.Filepath)
		}

		content, ok := buffers[sub.Filepath]
		if !ok {
			data, err := os.ReadFile(sub.Filepath)
			if err != nil {
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

		buffers[sub.Filepath] = strings.Replace(content, sub.Search, sub.Replace, 1)
	}

	return buffers, nil
}
