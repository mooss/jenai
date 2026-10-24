package main

import (
	"fmt"
	"os"

	"github.com/mooss/bagend/go/codec"
	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/bagend/go/sysu"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/highlight"
	"github.com/mooss/jenai/go/utils"

	"gopkg.in/yaml.v3"
)

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

func noerr[T any](res T, err error) T {
	if err != nil {
		fatal(err)
	}
	return res
}

func noerr0(err error) { noerr(0, err) }

func pretty(data any) string {
	formatter := utils.ToYAML[any]
	if config.Global.OutputFormat == "json" {
		formatter = codec.ToIndentedJSON[any]
	}

	formatted := noerr(formatter(data))
	formatted = lie.EnsureFinal(formatted, '\n')
	return string(formatted)
}

func prettyPrint(data any) {
	formatted := pretty(data)
	if sysu.TerminalColorEnabled() {
		formatted = highlight.Highlight(config.Global.OutputFormat, formatted)
	}

	fmt.Print(formatted)
}

func structToMap(data any) (map[string]any, error) {
	yamlbytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	err = yaml.Unmarshal(yamlbytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
