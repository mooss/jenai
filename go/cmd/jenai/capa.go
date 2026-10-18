package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/utils"
)

func runCapa(cfg *config.Jenai, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: jenai capa MODEL")
	}

	spec, err := modelSpec(&config.Jenai{Model: args[0]})
	if err != nil {
		return err
	}

	if spec.Provider != "openrouter" {
		return fmt.Errorf("unsupported provider %s (only openrouter is supported)", spec.Provider)
	}

	capabilities, err := openrouter.FetchCapabilities(
		context.Background(), spec, cfg.APIKeys[spec.Provider],
	)
	if err != nil {
		return err
	}

	model, err := utils.FromJSON[any](capabilities.Model)
	if err != nil {
		return err
	}

	endpoints, err := utils.FromJSON[any](capabilities.Endpoints)
	if err != nil {
		return err
	}

	fmt.Println(pretty(map[string]any{
		"model":     model,
		"endpoints": endpoints,
	}))

	return nil
}
