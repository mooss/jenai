package renderer

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mooss/bagend/go/sysu"
	"github.com/mooss/jenai/go/highlight"
)

// MarkdownRenderer streams completion content as syntax-highlighted markdown to a destination
// writer, with optional rendering of model reasoning.
type MarkdownRenderer struct {
	destination     io.Writer
	endsWithNewline bool
	output          *highlight.MarkdownStreamer
	reasoningOutput *highlight.MarkdownStreamer
	think           thinkState
}

// NewMarkdownRenderer returns a renderer writing to destination.
// If thoughtless is true, reasoning chunks are silently dropped.
func NewMarkdownRenderer(destination io.Writer, thoughtless bool) *MarkdownRenderer {
	color := sysu.TerminalColorEnabled()
	startingState := pending
	if thoughtless {
		startingState = over
	}

	return &MarkdownRenderer{
		destination:     destination,
		endsWithNewline: true,
		output:          highlight.NewMarkdownStreamer(destination, color),
		reasoningOutput: highlight.NewMarkdownStreamer(destination, color),
		think:           startingState,
	}
}

// ReadReasoning renders a chunk of model reasoning.
func (mdr *MarkdownRenderer) ReadReasoning(s string) error {
	//nolint:exhaustive
	switch mdr.think {
	case pending:
		mdr.think = thinking
		tag := mdr.reasoningOutput.ColorizeTag("<think>\n")
		_, err := fmt.Fprint(mdr.destination, tag)
		if err != nil {
			return err
		}
		fallthrough
	case thinking:
		return mdr.reasoningOutput.Write(s)
	default:
		// We'll just assume that requests are correctly formed and no reasoning happens after
		// content starts.
	}
	return nil
}

// ReadContent renders a chunk of completion content.
func (mdr *MarkdownRenderer) ReadContent(s string) error {
	//nolint:exhaustive
	switch mdr.think {
	case pending:
		mdr.think = over
	case thinking:
		mdr.think = over
		if err := mdr.reasoningOutput.Flush(); err != nil {
			return err
		}
		_, err := fmt.Fprint(mdr.destination, mdr.reasoningOutput.ColorizeTag("\n</think>\n"))
		if err != nil {
			return err
		}
	default:
	}

	mdr.endsWithNewline = strings.HasSuffix(s, "\n")
	return mdr.output.Write(s)
}

// UnclosedReasoning reports whether reasoning was started but never followed by content.
func (mdr *MarkdownRenderer) UnclosedReasoning() bool {
	return mdr.think == thinking
}

// Finish flushes the output, appending a trailing newline if missing.
func (mdr *MarkdownRenderer) Finish() error {
	var err error
	if !mdr.endsWithNewline {
		err = mdr.ReadContent("\n")
	}
	return errors.Join(err, mdr.output.Flush())
}
