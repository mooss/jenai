package config

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

///////////////////
// Context class //

// Context handles inclusion of files and directories as context.
type Context struct {
	Files       []string `yaml:"files,omitempty" json:"files,omitempty"`
	Dirs        []string `yaml:"dirs,omitempty" json:"dirs,omitempty"`
	Above       bool     `yaml:"above,omitempty" json:"above,omitempty"`
	LineNumbers bool     `yaml:"linenumbers,omitempty" json:"linenumbers,omitempty"`
	Title       string   `yaml:"-" json:"-"`
}

func (c *Context) Empty() bool { return len(c.Files) == 0 && len(c.Dirs) == 0 }

// Build returns the context string from included paths and the included paths.
// Empty should be checked before calling this because the "additional context" header is always
// present,
func (c *Context) Build() (string, []string, error) {
	if c.Empty() {
		return "", nil, nil
	}

	var buf bytes.Buffer
	paths := []string{}
	title := c.Title
	if title == "" {
		title = "Additional context (files)"
	}
	buf.WriteString("## ")
	buf.WriteString(title)

	for path, err := range c.allPaths {
		if err != nil {
			return "", nil, c.wrap(err)
		}

		if err := fileContent(&buf, path, c.LineNumbers); err != nil {
			return "", nil, c.wrap(err)
		}

		paths = append(paths, path)
	}

	return buf.String(), paths, nil
}

// allPaths returns an iterator on all paths contained in the context.
func (c *Context) allPaths(yield func(string, error) bool) {
	for _, path := range c.Files {
		if !yield(path, nil) {
			return
		}
	}

	for _, dir := range c.Dirs {
		if !yieldDirPaths(dir, yield) {
			return
		}
	}
}

func yieldDirPaths(dir string, yield func(string, error) bool) bool {
	for path, err := range iterFiles(dir) {
		if !yield(path, err) {
			return false
		}
	}
	return true
}

func (*Context) wrap(err error) error {
	return fmt.Errorf("failed to build file context: %w", err)
}

///////////////////////
// Utility functions //

// fileContent reads a file or URL and fills the buffer with its content, along with a header.
//
//nolint:revive
func fileContent(buf *bytes.Buffer, path string, linum bool) error {
	reader, err := readContent(path)
	if err != nil {
		return fmt.Errorf("error reading content from %s: %w", path, err)
	}
	defer reader.Close()

	buf.WriteString("\n\n====> START OF ")
	buf.WriteString(path)
	buf.WriteString(" <====\n\n")

	copyfun := func(dst *bytes.Buffer, src io.Reader) error {
		_, err := io.Copy(dst, src)
		return err
	}

	if linum {
		copyfun = copyWithLinum
	}

	if err := copyfun(buf, reader); err != nil {
		return fmt.Errorf("failed to copy %s content to buffer: %w", path, err)
	}

	buf.WriteString("\n\n====> END OF ")
	buf.WriteString(path)
	buf.WriteString(" <====")
	return nil
}

func copyWithLinum(buf *bytes.Buffer, src io.Reader) error {
	scanner := bufio.NewScanner(src)
	lineNumber := 1
	for scanner.Scan() {
		buf.WriteString(strconv.Itoa(lineNumber))
		buf.WriteString(": ")
		buf.Write(scanner.Bytes())
		buf.WriteRune('\n')
		lineNumber++
	}

	return scanner.Err()
}

// readContent returns an io.Reader for the given path, which can be a local file or a URL.
// It also returns a cleanup function to be called when the reader is no longer needed.
func readContent(path string) (io.ReadCloser, error) {
	if isURL(path) {
		return downloadContent(path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// isURL checks if the given string is a valid URL.
func isURL(path string) bool {
	u, err := url.Parse(path)
	return err == nil && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https")
}

// downloadContent downloads content from the given URL.
// The caller must close the returned reader.
func downloadContent(fileURL string) (body io.ReadCloser, err error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download URL %s: %w", fileURL, err)
	}

	client := *http.DefaultClient
	client.Timeout = 5 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download URL %s: %w", fileURL, err)
	}
	defer func() {
		if err != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		msg := "failed to download URL %s, status code: %d"
		return nil, fmt.Errorf(msg, fileURL, resp.StatusCode)
	}

	return resp.Body, nil
}

// iterFiles returns an iterator that yields all regular file paths under the given root directory.
//
//nolint:revive
func iterFiles(root string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				if !yield(path, err) {
					// The caller decides whether to stop.
					return err
				}

				return nil
			}

			if info.IsDir() {
				return nil
			}

			if !yield(path, nil) {
				return errors.New("IGNORED")
			}

			return nil
		})
	}
}
