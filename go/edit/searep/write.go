package searep

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// UpdatedFiles maps file paths to their new content.
type UpdatedFiles map[string]string

// Write writes all updated files to disk.
func (upd UpdatedFiles) Write() error {
	staged := make([]stagedFile, 0, len(upd))
	defer func() {
		for _, file := range staged {
			_ = os.Remove(file.temp)
		}
	}()

	var errs []error
	for path, content := range upd {
		file, err := stageEditableFile(path, content)
		if err != nil {
			log.Printf("stage editable file %q: %v", path, err)
			errs = append(errs, err)
			if file.temp != "" {
				_ = os.Remove(file.temp)
			}
			continue
		}
		staged = append(staged, file)
	}

	if err := replaceStagedFiles(staged); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// stagedFile is a file staged in a temporary location, pending rename to its final path.
type stagedFile struct {
	path string
	temp string
}

// stageEditableFile writes content to a temporary file in the same directory as path, ready to
// create or replace it.
func stageEditableFile(path, content string) (stagedFile, error) {
	mode, err := editableFileMode(path)
	if err != nil {
		return stagedFile{}, err
	}

	temp, err := os.CreateTemp(filepath.Dir(path), ".jenai-edit-*")
	if err != nil {
		return stagedFile{}, fmt.Errorf("stage editable file %q: %w", path, err)
	}

	file := stagedFile{path: path, temp: temp.Name()}
	if err := temp.Chmod(mode); err != nil {
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

// editableFileMode returns the file mode of path, or a default mode if the file does not exist yet.
func editableFileMode(path string) (fs.FileMode, error) {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return 0o666, nil
	}
	if err != nil {
		return 0, fmt.Errorf("stat editable file %q: %w", path, err)
	}

	return info.Mode(), nil
}

// replaceStagedFiles renames all staged temporary files to their final paths, replacing the
// original files.
func replaceStagedFiles(staged []stagedFile) error {
	for _, file := range staged {
		if err := os.Rename(file.temp, file.path); err != nil {
			return fmt.Errorf("replace editable file %q: %w", file.path, err)
		}
	}
	return nil
}
