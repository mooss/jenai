package workspace

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mooss/bagend/go/obj/oset"
	"github.com/mooss/bagend/go/pth"
)

type FileSection struct {
	Files []string `yaml:"files,omitempty" json:"files,omitempty"`
	Dirs  []string `yaml:"dirs,omitempty" json:"dirs,omitempty"`
	Title string   `yaml:"-" json:"-"`
}

func (s *FileSection) Empty() bool {
	return len(s.Files) == 0 && len(s.Dirs) == 0
}

func (s *FileSection) titleOr(defaultTitle string) string {
	if s.Title != "" {
		return s.Title
	}
	return defaultTitle
}

// readablePaths expands files and dirs and resolves them to absolute paths.
func (s *FileSection) readablePaths() ([]string, error) {
	raw, err := s.listAll()
	if err != nil {
		return nil, err
	}
	resolved := make([]string, 0, len(raw))
	for _, path := range raw {
		res, err := pth.ResolveRel(path)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, res)
	}
	return dedupOrdered(resolved), nil
}

// writablePaths expands files and dirs and splits them by existence.
// ResolveRel requires existence, so missing paths use resolveCreatablePath.
//
//revive:disable:cognitive-complexity
func (s *FileSection) writablePaths() (writable, creatable []string, err error) {
	raw, err := s.listAll()
	if err != nil {
		return nil, nil, err
	}
	seen := make(map[string]bool)
	for _, rawPath := range raw {
		if strings.TrimSpace(rawPath) == "" {
			return nil, nil, errors.New("writable files cannot be empty")
		}
		exists := pth.Exists(rawPath)
		var path string
		if exists {
			path, err = pth.ResolveRel(rawPath)
		} else {
			path, err = resolveCreatablePath(rawPath)
		}
		if err != nil {
			return nil, nil, err
		}
		if seen[path] {
			continue
		}
		seen[path] = true
		if exists {
			writable = append(writable, path)
		} else {
			creatable = append(creatable, path)
		}
	}
	return writable, creatable, nil
}

// listAll returns configured files plus every file under configured dirs.
func (s *FileSection) listAll() ([]string, error) {
	out := slices.Clone(s.Files)
	for _, dir := range s.Dirs {
		files, err := collectDirFiles(dir)
		if err != nil {
			return nil, err
		}
		out = append(out, files...)
	}
	return out, nil
}

/////////////////////
// Local utilities //

// dedupOrdered returns distinct paths in first-seen order.
func dedupOrdered(paths []string) []string {
	seen := oset.Empty[string]()
	return slices.Collect(seen.AddUniq(slices.Values(paths)))
}

// resolveCreatablePath resolves the nearest existing ancestor and appends the missing components.
func resolveCreatablePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	parent := abs
	var missing []string
	for !pth.Exists(parent) {
		next := filepath.Dir(parent)
		if next == parent {
			return "", fmt.Errorf("no existing ancestor for %s", path)
		}
		missing = append(missing, filepath.Base(parent))
		parent = next
	}
	resolved, err := pth.ResolveRel(parent)
	if err != nil {
		return "", err
	}
	for _, el := range slices.Backward(missing) {
		resolved = filepath.Join(resolved, el)
	}
	return resolved, nil
}
