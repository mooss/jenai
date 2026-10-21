package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mooss/bagend/go/fun/eager/lie"
)

var (
	seekMarker       = "[[[[[[[ SEEK"
	separatorMarker  = "|||||||"
	substituteMarker = "]]]]]]] SUBSTITUTE"
)

type substitution struct {
	File       string
	Seek       string
	Substitute string
}

/////////////
// Parsing //

// parseSubstitutions parses a response string into substitutions.
// This is very procedural but splitting it would make it harder to follow.
//
//revive:disable:cognitive-complexity,cyclomatic
func parseSubstitutions(response string) ([]substitution, error) {
	lines := strings.Split(response, "\n")
	var result []substitution

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
			file == seekMarker || file == separatorMarker || file == substituteMarker {
			return nil, fmt.Errorf("invalid substitution file path on line %d", i+1)
		}
		i++

		// Code block fence.
		if i >= len(lines) || !strings.HasPrefix(lines[i], "```") {
			return nil, fmt.Errorf("expected opening code fence after file path %q", file)
		}
		i++
		if i >= len(lines) || lines[i] != seekMarker {
			return nil, fmt.Errorf("expected seek marker for %q", file)
		}
		i++

		// Seek separator.
		seekLines, next, ok := collectUntil(lines, i, separatorMarker)
		if !ok {
			return nil, fmt.Errorf("missing separator marker for %q", file)
		}
		i = next + 1

		// Seek end of substitute section.
		substituteLines, next, ok := collectUntil(lines, i, substituteMarker)
		if !ok {
			return nil, fmt.Errorf("missing substitute marker for %q", file)
		}
		i = next + 1

		// Closing code block fence.
		if i >= len(lines) || !strings.HasPrefix(lines[i], "```") {
			return nil, fmt.Errorf("expected closing code fence for %q", file)
		}
		i++

		result = append(result, substitution{
			File:       file,
			Seek:       joinBlockLines(seekLines),
			Substitute: joinBlockLines(substituteLines),
		})
	}

	if len(result) == 0 {
		return nil, errors.New("response contains no complete seek/substitute blocks")
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

////////////////
// Validating //

func validateSubstitutions(
	subs []substitution, writeFiles []string,
) (map[string]string, error) {
	buffers := make(map[string]string)
	allowed := make(map[string]struct{}, len(writeFiles))
	for _, file := range writeFiles {
		resolved, err := resolveWritePath(file)
		if err != nil {
			return nil, err
		}
		allowed[resolved] = struct{}{}
	}

	for _, sub := range subs {
		resolved, err := resolveWritePath(sub.File)
		if err != nil {
			return nil, err
		}
		_, ok := allowed[resolved]
		if !ok {
			msg := "refusing attempt to write to file %q (resolved to %q, not in --write)"
			return nil, fmt.Errorf(msg, sub.File, resolved)
		}

		content, ok := buffers[resolved]
		if !ok {
			data, err := os.ReadFile(resolved)
			if err != nil {
				return nil, fmt.Errorf("when reading editable file %q: %w", resolved, err)
			}
			// Files must end in newlines otherwise seek/substitute could break on the last line.
			content = string(lie.EnsureFinal(data, '\n'))
		}

		if sub.Seek == "" {
			return nil, fmt.Errorf("seek content for %q must not be empty", sub.File)
		}
		if count := strings.Count(content, sub.Seek); count != 1 {
			msg := "seek content for %q occurs %d times; expected exactly once"
			return nil, fmt.Errorf(msg, sub.File, count)
		}

		buffers[resolved] = strings.Replace(content, sub.Seek, sub.Substitute, 1)
	}

	return buffers, nil
}

func writeBuffers(buffers map[string]string) error {
	staged := make([]stagedFile, 0, len(buffers))
	defer func() {
		for _, file := range staged {
			_ = os.Remove(file.temp)
		}
	}()

	for path, content := range buffers {
		file, err := stageEditableFile(path, content)
		staged = append(staged, file)
		if err != nil {
			return err
		}
	}

	return replaceStagedFiles(staged)
}

type stagedFile struct {
	path string
	temp string
}

func stageEditableFile(path, content string) (stagedFile, error) {
	resolved, err := resolveWritePath(path)
	if err != nil {
		return stagedFile{}, err
	}
	path = resolved

	info, err := os.Stat(path)
	if err != nil {
		return stagedFile{}, fmt.Errorf("stat editable file %q: %w", path, err)
	}

	temp, err := os.CreateTemp(filepath.Dir(path), ".jenai-edit-*")
	if err != nil {
		return stagedFile{}, fmt.Errorf("stage editable file %q: %w", path, err)
	}

	file := stagedFile{path: path, temp: temp.Name()}
	if err := temp.Chmod(info.Mode()); err != nil {
		_ = temp.Close()
		return file, fmt.Errorf("set permissions for %q: %w", path, err)
	}

	if _, err := temp.WriteString(content); err != nil {
		_ = temp.Close()
		return file, fmt.Errorf("write staged file %q: %w", path, err)
	}

	if err := temp.Close(); err != nil {
		return file, fmt.Errorf("close staged file %q: %w", path, err)
	}

	return file, nil
}

func replaceStagedFiles(staged []stagedFile) error {
	for _, file := range staged {
		if err := os.Rename(file.temp, file.path); err != nil {
			return fmt.Errorf("replace editable file %q: %w", file.path, err)
		}
	}
	return nil
}
