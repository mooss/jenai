package main

import "github.com/mooss/jenai/go/config"

func runEdit(qry *config.Query, paths []string, force bool) error {
	prepare := func() error {
		return validateWriteFiles(paths, force)
	}
	mkprompt := func() (config.Prompt, error) {
		return buildEditPrompt(qry, paths)
	}
	onResponse := func(response string) error {
		substitutions, err := parseSubstitutions(response)
		if err != nil {
			return err
		}
		buffers, err := validateSubstitutions(substitutions, paths)
		if err != nil {
			return err
		}
		return writeBuffers(buffers)
	}

	return runWithHooks(qry, mkprompt, prepare, onResponse)
}
