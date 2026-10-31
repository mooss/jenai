package renderer

import (
	"strings"

	"github.com/mooss/jenai/go/utils"
)

// TeeRenderer renders completion to optional files.
// It cannot fail, though TeeFile can sometimes warn if it cannot write to its output.
type TeeRenderer struct {
	full            *utils.TeeFile
	thoughtless     *utils.TeeFile
	think           thinkState
	endsWithNewline bool
}

func NewTeeRenderer(teeFile, thoughtlessFile *utils.TeeFile) *TeeRenderer {
	return &TeeRenderer{
		full:            teeFile,
		thoughtless:     thoughtlessFile,
		think:           pending,
		endsWithNewline: true,
	}
}

// ReadReasoning records reasoning to the full tee file only.
func (ter *TeeRenderer) ReadReasoning(s string) error {
	//nolint:exhaustive
	switch ter.think {
	case pending:
		ter.full.Write("<think>\n")
		ter.think = thinking
	case over:
		// We'll just assume that requests are correctly formed and no reasoning happens after
		// content starts.
		return nil
	default:
	}

	ter.full.Write(s)
	return nil
}

// ReadContent records content to both tee files.
func (ter *TeeRenderer) ReadContent(s string) error {
	//nolint:exhaustive
	switch ter.think {
	case pending:
		ter.think = over
	case thinking:
		ter.think = over
		ter.full.Write("\n</think>\n")
	default:
	}

	ter.endsWithNewline = strings.HasSuffix(s, "\n")
	ter.full.Write(s)
	ter.thoughtless.Write(s)
	return nil
}

// Finish closes both tee files, appending a trailing newline if missing.
func (ter *TeeRenderer) Finish() error {
	if !ter.endsWithNewline {
		ter.ReadContent("\n")
	}
	ter.full.Finish()
	ter.thoughtless.Finish()
	return nil
}
