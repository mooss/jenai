package main

import (
	"bytes"
	"fmt"
	"os"

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
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	noerr0(enc.Encode(data))
	noerr0(enc.Close())
	return buf.String()
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
