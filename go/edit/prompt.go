package edit

import (
	"fmt"

	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/edit/searep"
)

const editPromptText = `Edit the requested files according to the user's instructions.

Return only one or more search/replace blocks in the exact format below.
Do not include anything but those search/replace blocks.
Use only file paths from the ` + "`## Editable files`" + ` section.

FILE PATH
` + "```language" + `
%s
SEARCH CONTENT
%s
REPLACEMENT CONTENT
%s
` + "```" + `

For instance:

go/flag/parser.go
` + "```go" + `
%s
			dest = &par.Positional
%s
			dest = &par.positional
%s
` + "```" + `

The search content must match the file verbatim and appear exactly once.
Include more surrounding context if needed to make it unique.
Preserve the file's indentation and style.
The replacement content is the exact text that will replace the search content.
Make sure to close the code block fence.`

func BuildPrompt(qry *config.Query) (config.Prompt, error) {
	qry.Context.Readable.Title = "Read-only files"
	prompt, err := qry.BuildOneshotPrompt()
	if err != nil {
		return config.Prompt{}, err
	}

	instructions := fmt.Sprintf(
		editPromptText,
		searep.SearchMarker, searep.SeparatorMarker, searep.ReplaceMarker,
		searep.SearchMarker, searep.SeparatorMarker, searep.ReplaceMarker,
	)
	prompt.Primary = instructions

	return prompt, nil
}
