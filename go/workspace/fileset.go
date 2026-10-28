package workspace

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/mooss/bagend/go/obj/oset"
)

// Fileset handles the interaction with files and directories.
type Fileset struct {
	Readable    FileSection `yaml:"readable" json:"readable"`
	Writable    FileSection `yaml:"writable" json:"writable"`
	Above       bool        `yaml:"above,omitempty" json:"above,omitempty"`
	LineNumbers bool        `yaml:"linenumbers,omitempty" json:"linenumbers,omitempty"`

	resolved *resolvedCache
}

func (c *Fileset) Empty() bool {
	return c.Readable.Empty() && c.Writable.Empty()
}

// Resolve expands files and dirs once and splits writable paths by existence.
func (c *Fileset) Resolve() (ResolvedPaths, error) {
	if c.resolved == nil {
		res, err := c.resolveImpl()
		c.resolved = &resolvedCache{res, err}
	}

	return c.resolved.paths, c.resolved.err
}

func (c *Fileset) resolveImpl() (ResolvedPaths, error) {
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
func (c *Fileset) Build() (string, []string, error) {
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

func (c *Fileset) writeSection(
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
func (c *Fileset) writeSectionHeader(buf *bytes.Buffer, title string, linum bool, first bool) {
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

func (*Fileset) wrap(err error) error {
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
// fileContent //

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
