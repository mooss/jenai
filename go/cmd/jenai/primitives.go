package main

import (
	"fmt"
	"os"

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
	return string(noerr(utils.ToYAML(data)))
}

func prettyPrint(data any) {
	fmt.Print(pretty(data))
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
