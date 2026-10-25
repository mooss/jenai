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
	"strings"
	"time"
)

///////////////////
// Context class //

// Context handles inclusion of files and directories as context.
type Context struct {
	Readable    FileSection `yaml:"readable" json:"readable"`
	Writable    FileSection `yaml:"writable" json:"writable"`
	Above       bool        `yaml:"above,omitempty" json:"above,omitempty"`
	LineNumbers bool        `yaml:"linenumbers,omitempty" json:"linenumbers,omitempty"`
}

func (c *Context) Empty() bool {
	return c.Readable.Empty() && c.Writable.Empty()
}

// Build returns the context string from included paths and the included paths.
// Empty should be checked before calling this because the "additional context" header is always
// present,
func (c *Context) Build() (string, []string, error) {
	if c.Empty() {
		return "", nil, nil
	}

	var buf bytes.Buffer
	paths := []string{}

	if !c.Readable.Empty() {
		title := c.Readable.titleOr("Additional context (files)")
		err := c.writeSection(&buf, &paths, title, c.Readable.AllPaths, c.LineNumbers)
		if err != nil {
			return "", nil, err
		}
	}

	if !c.Writable.Empty() {
		title := c.Writable.titleOr("Editable files")
		if err := c.writeSection(&buf, &paths, title, c.Writable.AllPaths, false); err != nil {
			return "", nil, err
		}
	}

	return buf.String(), paths, nil
}

func (c *Context) writeSection(
	buf *bytes.Buffer, paths *[]string, title string, allPaths iter.Seq2[string, error], linum bool,
) error {
	c.writeSectionHeader(buf, title, linum)

	for path, err := range allPaths {
		if err != nil {
			return c.wrap(err)
		}
		if err := fileContent(buf, path, linum); err != nil {
			return c.wrap(err)
		}
		*paths = append(*paths, path)
	}
	return nil
}

//revive:disable:line-length-limit,flag-parameter
func (*Context) writeSectionHeader(buf *bytes.Buffer, title string, linum bool) {
	buf.WriteString("## ")
	buf.WriteString(title)
	buf.WriteString("\n\nFiles are provided in Markdown code blocks, with a `### Start of $FILENAME` section above and a `> End of $FILENAME` quote below.")
	buf.WriteString("\nThe Markdown code block fences (triple backticks) are not part of the files.")
	if linum {
		buf.WriteString("\nEach line of a file is prefixed with the line number in the format `line: `.")
	}
	buf.WriteString("\n")
}

func (*Context) wrap(err error) error {
	return fmt.Errorf("failed to build file context: %w", err)
}

/////////////////
// FileSection //

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

// AllPaths returns an iterator on all paths contained in the section.
func (s *FileSection) AllPaths(yield func(string, error) bool) {
	for _, path := range s.Files {
		if !yield(path, nil) {
			return
		}
	}

	for _, dir := range s.Dirs {
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

///////////////////////
// Utility functions //

// fileContent reads a file or URL and fills the buffer with its content, along with a header.
func fileContent(buf *bytes.Buffer, path string, linum bool) error {
	reader, err := readContent(path)
	if err != nil {
		return fmt.Errorf("error reading content from %s: %w", path, err)
	}
	defer reader.Close()

	buf.WriteString("\n### Start of ")
	buf.WriteString(path)
	buf.WriteString("\n```")
	buf.WriteString(strings.TrimPrefix(filepath.Ext(path), "."))
	buf.WriteRune('\n')

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

	buf.WriteString("\n```\n> End of ")
	buf.WriteString(path)
	buf.WriteString("\n")
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
	client.Timeout = 2 * time.Second

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
//revive:disable:cognitive-complexity
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
