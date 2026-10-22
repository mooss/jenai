package main

import (
	"fmt"

	"github.com/mooss/jenai/go/config"
)

const editPromptText = `Edit the requested files according to the user's instructions.

Return only one or more seek/substitute blocks in the exact format below.
Do not include anything but those seek/substitute blocks.
Use only file paths from the EDITABLE FILES section.

FILE-PATH
` + "```language" + `
%s
SEEK CONTENT
%s
REPLACEMENT CONTENT
%s
` + "```" + `

The seek content must match the file verbatim and appear exactly once.
Include more surrounding context if needed to make it unique.
Preserve the file's indentation and style.
The replacement content is the exact text that will replace the seek content.`

func buildEditPrompt(qry *config.Query, writeFiles []string) (config.Prompt, error) {
	prompt, err := qry.BuildOneshotPrompt()
	if err != nil {
		return config.Prompt{}, err
	}

	instructions := fmt.Sprintf(editPromptText, seekMarker, separatorMarker, substituteMarker)
	prompt.Primary = instructions
	context := config.Context{Files: writeFiles, Title: "EDITABLE FILES"}

	content, _, err := context.Build()
	if err != nil {
		return config.Prompt{}, err
	}

	if prompt.Context == "" {
		prompt.Context = content
	} else {
		prompt.Context += "\n\n" + content
	}

	return prompt, nil
}
