package renderer

import (
	"errors"
	"os"
	"strings"

	"github.com/mooss/jenai/go/utils"
)

// TerminalRenderer renders a completion to the terminal while recording the response and teeing the
// output to optional files.
type TerminalRenderer struct {
	markdown *MarkdownRenderer
	tee      *TeeRenderer
	response strings.Builder
}

// NewTerminalRenderer renders to stdout, teeing to the given files.
// If thoughtless is true, reasoning is not rendered to stdout.
func NewTerminalRenderer(
	teeFile, thoughtlessFile *utils.TeeFile, thoughtless bool,
) *TerminalRenderer {
	return &TerminalRenderer{
		markdown: NewMarkdownRenderer(os.Stdout, thoughtless),
		tee:      NewTeeRenderer(teeFile, thoughtlessFile),
	}
}

// ReadReasoning renders a chunk of model reasoning.
func (r *TerminalRenderer) ReadReasoning(s string) error {
	r.tee.ReadReasoning(s) //nolint:errcheck // TeeRenderer never fails.
	return r.markdown.ReadReasoning(s)
}

// ReadContent renders a chunk of completion content.
func (r *TerminalRenderer) ReadContent(s string) error {
	r.response.WriteString(s)
	r.tee.ReadContent(s) //nolint:errcheck // TeeRenderer never fails.
	return r.markdown.ReadContent(s)
}

// Response returns the full completion content accumulated.
func (r *TerminalRenderer) Response() string {
	return r.response.String()
}

// UnclosedReasoning reports whether reasoning was started but never followed by content.
func (r *TerminalRenderer) UnclosedReasoning() bool {
	return r.markdown.UnclosedReasoning()
}

// Finish flushes the rendered output and closes the tee files.
func (r *TerminalRenderer) Finish() error {
	return errors.Join(r.markdown.Finish(), r.tee.Finish())
}
