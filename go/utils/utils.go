// Package utils provides utility that don't warrant their dedicated package for now.
// This is mostly a purgatory for utilities meant to eventually go to bagend.
package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// FromYAML unmarshalls YAML data to the given type.
func FromYAML[T any](data []byte) (T, error) {
	var res T
	if err := yaml.Unmarshal(data, &res); err != nil {
		var zero T
		return zero, err
	}

	return res, nil
}

// FromYAMLFile reads YAML data from a file and unmarshalls it to the given type.
func FromYAMLFile[T any](datapath string) (T, error) {
	data, err := os.ReadFile(datapath)
	if err != nil {
		var zero T
		return zero, err
	}

	return FromYAML[T](data)
}

// ToYAML marshalls arbitrary data to YAML.
func ToYAML[T any](source T) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	if err := enc.Encode(source); err != nil {
		return nil, err
	}

	if err := enc.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Wrapf wraps a non-nil err and prefixes it with a formatted string.
// It is meant to be a shortcut for wrapping a 2-value error result.
func Wrapf[T any](res T, err error, format string, args ...any) (T, error) {
	if err != nil {
		return res, fmt.Errorf(format+": %w", append(args, err)...)
	}

	return res, nil
}

// Wrap wraps a non-nil *err and prefixes it with a formatted string.
// It is meant to be called at the beginning of a function using defer.
func Wrap(err *error, format string, args ...any) {
	if err == nil || *err == nil {
		return
	}

	*err = fmt.Errorf(format+": %w", append(args, *err)...)
}

// OnceErr returns a function that calls mk exactly once and caches the result.
// Subsequent calls will return the cached value.
// The returned function is safe for concurrent use.
func OnceErr[T any](mk func() (T, error)) func() (T, error) {
	type cache struct {
		value T
		err   error
	}
	get := sync.OnceValue(func() cache {
		value, err := mk()
		return cache{value, err}
	})
	return func() (T, error) {
		res := get()
		return res.value, res.err
	}
}

// TeeFile is a buffered file writer that warns to stderr but never fails.
// Write and Finish are safe to call on nil and after failure/close.
// Not meant for concurrent use.
// Valid to use when nil (no-op).
type TeeFile struct {
	file   *os.File
	writer *bufio.Writer
	path   string
	closed bool
}

// OpenTeeFile opens path for writing with a buffered writer.
// It returns nil when path is empty or the file cannot be opened.
func OpenTeeFile(path string) *TeeFile {
	if path == "" {
		return nil
	}

	full := path
	if abs, err := filepath.Abs(path); err == nil {
		full = abs
	}

	file, err := os.OpenFile(full, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		teeWarn(path, err)
		return nil
	}

	return &TeeFile{file: file, writer: bufio.NewWriter(file), path: path}
}

// IsOpen reports whether the file is ready for writing.
func (t *TeeFile) IsOpen() bool {
	return t != nil && !t.closed
}

// Write buffers s, disabling further writes on error.
func (t *TeeFile) Write(s string) {
	if !t.IsOpen() {
		return
	}
	if _, err := t.writer.WriteString(s); err != nil {
		teeWarn(t.path, err)
		t.closeNow()
	}
}

// Finish flushes the buffer and closes the file.
func (t *TeeFile) Finish() {
	if !t.IsOpen() {
		return
	}
	if err := t.writer.Flush(); err != nil {
		teeWarn(t.path, err)
	}
	t.closeNow()
}

func teeWarn(path string, err error) {
	fmt.Fprintf(os.Stderr, "can't tee to %s (will proceed nonetheless): %s\n", path, err)
}

func (t *TeeFile) closeNow() {
	if err := t.file.Close(); err != nil {
		teeWarn(t.path, err)
	}
	t.closed = true
}
