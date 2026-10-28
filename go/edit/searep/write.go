package searep

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Buffers map[string]string

func (buf Buffers) Write() error {
	staged := make([]stagedFile, 0, len(buf))
	defer func() {
		for _, file := range staged {
			_ = os.Remove(file.temp)
		}
	}()

	for path, content := range buf {
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

func replaceStagedFiles(staged []stagedFile) error {
	for _, file := range staged {
		if err := os.Rename(file.temp, file.path); err != nil {
			return fmt.Errorf("replace editable file %q: %w", file.path, err)
		}
	}
	return nil
}
