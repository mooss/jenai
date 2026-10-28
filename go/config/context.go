package config

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mooss/bagend/go/obj/oset"
	"github.com/mooss/bagend/go/pth"
)

///////////////////
// Context class //

// Context handles inclusion of files and directories as context.
type Context struct {
	Readable    FileSection `yaml:"readable" json:"readable"`
	Writable    FileSection `yaml:"writable" json:"writable"`
	Above       bool        `yaml:"above,omitempty" json:"above,omitempty"`
	LineNumbers bool        `yaml:"linenumbers,omitempty" json:"linenumbers,omitempty"`

	resolved *resolvedCache
}

func (c *Context) Empty() bool {
	return c.Readable.Empty() && c.Writable.Empty()
}

// ResolvedPaths holds deduplicated absolute paths, split by kind.
type ResolvedPaths struct {
	Readable  []string
	Writable  []string
	Creatable []string
}

// All returns every path, readable first.
func (p ResolvedPaths) All() []string {
	return slices.Concat(p.Readable, p.Writable, p.Creatable)
}

// Editable returns the writable and creatable paths.
func (p ResolvedPaths) Editable() []string {
	return slices.Concat(p.Writable, p.Creatable)
}

type resolvedCache struct {
	paths ResolvedPaths
	err   error
}

// Resolve expands files and dirs once and splits writable paths by existence.
func (c *Context) Resolve() (ResolvedPaths, error) {
	if c.resolved == nil {
		res, err := c.resolveImpl()
		c.resolved = &resolvedCache{res, err}
	}

	return c.resolved.paths, c.resolved.err
}

func (c *Context) resolveImpl() (ResolvedPaths, error) {
	readable, err := c.Readable.readablePaths()
	if err != nil {
		return ResolvedPaths{}, err
	}
	writable, creatable, err := c.Writable.writablePaths()
	if err != nil {
		return ResolvedPaths{}, err
	}
	return ResolvedPaths{
		Readable:  readable,
		Writable:  writable,
		Creatable: creatable,
	}, nil
}

// Build returns the context string from included paths and the included paths.
// Empty should be checked before calling this because the "additional context" header is always
// present,
func (c *Context) Build() (string, []string, error) {
	if c.Empty() {
		return "", nil, nil
	}

	paths, err := c.Resolve()
	if err != nil {
		return "", nil, c.wrap(err)
	}

	dupes := oset.New(paths.Readable...).Intersect(slices.Values(paths.Editable()))
	if !dupes.Empty() {
		msg := "files %v were provided as both readable and writable"
		return "", nil, c.wrap(fmt.Errorf(msg, slices.Sorted(dupes.Keys())))
	}

	var buf bytes.Buffer
	first := true

	if len(paths.Readable) > 0 {
		title := c.Readable.titleOr("Additional context (files)")
		if err := c.writeSection(&buf, title, paths.Readable, c.LineNumbers, first); err != nil {
			return "", nil, err
		}
		first = false
	}

	if len(paths.Writable) > 0 {
		title := c.Writable.titleOr("Editable files")
		if err := c.writeSection(&buf, title, paths.Writable, false, first); err != nil {
			return "", nil, err
		}
		first = false
	}

	if len(paths.Creatable) > 0 {
		writeCreatableSection(&buf, paths.Creatable, first)
	}

	return buf.String(), paths.All(), nil
}

//////////////
// Sections //

func (c *Context) writeSection(
	buf *bytes.Buffer, title string, paths []string, linum bool, first bool,
) error {
	c.writeSectionHeader(buf, title, linum, first)

	for _, path := range paths {
		if err := fileContent(buf, path, linum); err != nil {
			return c.wrap(err)
		}
	}
	return nil
}

//revive:disable:line-length-limit,flag-parameter
func (c *Context) writeSectionHeader(buf *bytes.Buffer, title string, linum bool, first bool) {
	if !first {
		buf.WriteRune('\n')
	}
	buf.WriteString("## ")
	buf.WriteString(title)
	if first {
		buf.WriteString("\n\nFiles are provided in Markdown code blocks, with a `### Start of $FILENAME` section above and a `> End of $FILENAME` quote below.")
		buf.WriteString("\nThe Markdown code block fences (triple backticks) just after the start and just before the end are not part of the files.")
		if linum {
			buf.WriteString("\nThe files in this section are prefixed with the line number in the format `line: `.")
		}
	} else if c.LineNumbers {
		// The first section contains line numbers so remind that this section in contrast does not.
		buf.WriteString("\n\nThe files in this section do not contain line numbers.")
	}
	buf.WriteString("\n")
}

func (*Context) wrap(err error) error {
	return fmt.Errorf("failed to build file context: %w", err)
}

// writeCreatableSection lists writable paths that do not exist yet and can be created.
func writeCreatableSection(buf *bytes.Buffer, paths []string, first bool) {
	if !first {
		buf.WriteRune('\n')
	}
	buf.WriteString("## Creatable files\n\n")
	buf.WriteString("The following writable files do not exist yet, but can be created with a search/replace block if needed.\n")
	for _, path := range paths {
		buf.WriteString("- ")
		buf.WriteString(path)
		buf.WriteRune('\n')
	}
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

	if linum {
		err = copyWithLinum(buf, reader)
	} else {
		_, err = io.Copy(buf, reader)
	}
	if err != nil {
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
func downloadContent(fileURL string) (io.ReadCloser, error) {
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

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		msg := "failed to download URL %s, status code: %d"
		return nil, fmt.Errorf(msg, fileURL, resp.StatusCode)
	}

	return resp.Body, nil
}

// collectDirFiles returns every regular file under root.
func collectDirFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
