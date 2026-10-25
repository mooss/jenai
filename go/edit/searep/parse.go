package searep

import (
	"errors"
	"fmt"
	"strings"
)

// Parse parses a response string into substitutions.
// This is very procedural but splitting it would make it harder to follow.
//
//revive:disable:cognitive-complexity,cyclomatic
func Parse(response string) ([]Block, error) {
	lines := strings.Split(response, "\n")
	var result []Block

	for i := 0; i < len(lines); {
		// Skip empty lines.
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i == len(lines) {
			break
		}

		// File path.
		file := lines[i]
		if file == "" || strings.TrimSpace(file) != file ||
			file == SearchMarker || file == SeparatorMarker || file == ReplaceMarker {
			return nil, fmt.Errorf("invalid substitution file path on line %d", i+1)
		}
		i++

		// Code block fence.
		if i >= len(lines) || !strings.HasPrefix(lines[i], "```") {
			return nil, fmt.Errorf("expected opening code fence after file path %q", file)
		}
		i++
		if i >= len(lines) || lines[i] != SearchMarker {
			return nil, fmt.Errorf("expected search marker for %q", file)
		}
		i++

		// Seek separator.
		searchLines, next, ok := collectUntil(lines, i, SeparatorMarker)
		if !ok {
			return nil, fmt.Errorf("missing separator marker for %q", file)
		}
		i = next + 1

		// Seek end of replacement section.
		replaceLines, next, ok := collectUntil(lines, i, ReplaceMarker)
		if !ok {
			return nil, fmt.Errorf("missing replacement marker for %q", file)
		}
		i = next + 1

		// Closing code block fence.
		if i >= len(lines) || !strings.HasPrefix(lines[i], "```") {
			return nil, fmt.Errorf("expected closing code fence for %q", file)
		}
		i++

		result = append(result, Block{
			File:    file,
			Search:  joinBlockLines(searchLines),
			Replace: joinBlockLines(replaceLines),
		})
	}

	if len(result) == 0 {
		return nil, errors.New("response contains no complete search/replace blocks")
	}

	return result, nil
}

func collectUntil(lines []string, start int, marker string) ([]string, int, bool) {
	for i := start; i < len(lines); i++ {
		if lines[i] == marker {
			return lines[start:i], i, true
		}
	}
	return nil, 0, false
}

func joinBlockLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}
