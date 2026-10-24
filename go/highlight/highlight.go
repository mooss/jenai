package highlight

import (
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var DefaultStyle = "pygments"

func Highlight(language, text string) string {
	lexer := lexers.Get(language)
	formatter := formatters.Get("terminal256")
	style := styles.Get(DefaultStyle)
	if lexer == nil || formatter == nil || style == nil {
		return text
	}

	tokens, err := lexer.Tokenise(nil, text)
	if err != nil {
		return text
	}

	var highlighted strings.Builder
	if err := formatter.Format(&highlighted, style, tokens); err != nil {
		return text
	}
	return highlighted.String()
}
