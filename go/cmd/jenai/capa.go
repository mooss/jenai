package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
)

func runCapa(cfg *config.Jenai, args []string, source string) error {
	spec, err := capaModelSpec(args)
	if err != nil {
		return err
	}

	ctx := context.Background()
	apiKey := cfg.APIKeys[spec.Provider]
	if source == "rawyaml" {
		return rawyaml(ctx, spec, apiKey)
	}

	endpointsResponse, err := openrouter.FetchCatalogRecord[openrouter.ListEndpointsResponse](
		ctx, modelURL(spec, "models")+"/endpoints", apiKey,
	)
	if err != nil {
		return fmt.Errorf("openrouter: fetch endpoint capabilities: %w", err)
	}

	if source == "endpoints" {
		endpoints := endpointsResponse.Data.Endpoints
		prettyPrint(endpoints)
		return nil
	}

	modelResponse, err := openrouter.FetchCatalogRecord[openrouter.GetModelResponse](
		ctx, modelURL(spec, "model"), apiKey,
	)
	if err != nil {
		return fmt.Errorf("openrouter: fetch model capabilities: %w", err)
	}

	endpointsResponse.Data.Endpoints = nil
	prettyPrint(map[string]any{
		"model":     modelResponse.Data,
		"endpoints": endpointsResponse,
	})
	return nil
}

func modelURL(spec models.Spec, resource string) string {
	return "https://openrouter.ai/api/v1/" + resource + "/" +
		url.PathEscape(spec.Author) + "/" + url.PathEscape(spec.Model)
}

func capaModelSpec(args []string) (models.Spec, error) {
	if len(args) != 1 {
		return models.Spec{}, errors.New("usage: jenai capa MODEL")
	}

	spec, err := modelSpec(&config.Jenai{Model: args[0]})
	if err != nil {
		return models.Spec{}, err
	}
	if spec.Provider != "openrouter" {
		msg := "unsupported provider %s (only openrouter is supported)"
		return models.Spec{}, fmt.Errorf(msg, spec.Provider)
	}

	return spec, nil
}

func rawyaml(ctx context.Context, spec models.Spec, apiKey string) error {
	endpoints, err := openrouter.FetchCatalogRecord[any](
		ctx, modelURL(spec, "models")+"/endpoints", apiKey,
	)
	if err != nil {
		return fmt.Errorf("openrouter: fetch endpoint capabilities: %w", err)
	}

	model, err := openrouter.FetchCatalogRecord[any](ctx, modelURL(spec, "model"), apiKey)
	if err != nil {
		return fmt.Errorf("openrouter: fetch model capabilities: %w", err)
	}

	prettyPrint(map[string]any{
		"model":     model,
		"endpoints": endpoints,
	})
	return nil
}
