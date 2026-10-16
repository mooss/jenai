package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func fatal(err error) {
	fmt.Println("Error:", err)
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
	return string(noerr(yaml.Marshal(data)))
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
