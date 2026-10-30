package main

import (
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/edit"
	"github.com/mooss/jenai/go/edit/searep"
)

func runEdit(qry *config.Query, force bool) error {
	resolved, err := qry.Context.Resolve()
	if err != nil {
		return err
	}
	paths := resolved.Editable()

	prepare := func() error {
		return edit.ValidateWritable(paths, force)
	}
	mkprompt := func() (config.Prompt, error) {
		return edit.BuildPrompt(qry)
	}
	onResponse := func(response string) error {
		blocks, err := searep.Parse(response)
		if err != nil {
			return err
		}
		buffers, err := searep.ApplySearchAndReplace(blocks, paths)
		if err != nil {
			return err
		}
		return buffers.Write()
	}

	return runWithHooks(qry, mkprompt, prepare, onResponse)
}
