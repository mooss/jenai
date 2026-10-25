// Package searep implement the search and replacement block protocol.
package searep

import (
	"fmt"
	"os"
	"strings"

	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/bagend/go/pth"
)

var (
	SearchMarker    = "<<<<<<< SEARCH"
	SeparatorMarker = "======="
	ReplaceMarker   = ">>>>>>> REPLACE"
)

type Block struct {
	File    string
	Search  string
	Replace string
}

func ApplyBlocks(subs []Block, writeFiles []string) (Buffers, error) {
	buffers := make(Buffers)
	allowed := make(map[string]struct{}, len(writeFiles))
	for _, file := range writeFiles {
		resolved, err := pth.Resolve(file)
		if err != nil {
			return nil, err
		}
		allowed[resolved] = struct{}{}
	}

	for _, sub := range subs {
		resolved, err := pth.Resolve(sub.File)
		if err != nil {
			return nil, err
		}
		_, ok := allowed[resolved]
		if !ok {
			msg := "refusing attempt to write to file %q (resolved to %q, not in writable files)"
			return nil, fmt.Errorf(msg, sub.File, resolved)
		}

		content, ok := buffers[resolved]
		if !ok {
			data, err := os.ReadFile(resolved)
			if err != nil {
				return nil, fmt.Errorf("when reading editable file %q: %w", resolved, err)
			}
			// Files must end in newlines otherwise search/replace could break on the last line.
			content = string(lie.EnsureFinal(data, '\n'))
		}

		if sub.Search == "" {
			return nil, fmt.Errorf("search content for %q must not be empty", sub.File)
		}
		if count := strings.Count(content, sub.Search); count != 1 {
			msg := "search content for %q occurs %d times; expected exactly once"
			return nil, fmt.Errorf(msg, sub.File, count)
		}

		buffers[resolved] = strings.Replace(content, sub.Search, sub.Replace, 1)
	}

	return buffers, nil
}
