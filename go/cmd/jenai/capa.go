package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
)

func runCapa(qry *config.Query, args []string, source string) error {
	spec, err := capaModelSpec(args)
	if err != nil {
		return err
	}

	ctx := context.Background()
	apiKey := qry.APIKeys[spec.Provider]

	switch source {
	case "rawyaml":
		return rawyaml(ctx, spec, apiKey)
	case "overview":
		return overview(ctx, spec, apiKey)
	}

	endpointsResponse, err := fetchRecord[openrouter.ListEndpointsResponse](
		ctx, modelURL(spec, "models")+"/endpoints", apiKey, "models endpoints",
	)
	if err != nil {
		return err
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

	spec, err := modelSpec(&config.Query{Model: args[0]})
	if err != nil {
		return models.Spec{}, err
	}
	if spec.Provider != "openrouter" {
		msg := "unsupported provider %s (only openrouter is supported)"
		return models.Spec{}, fmt.Errorf(msg, spec.Provider)
	}

	return spec, nil
}

func fetchRecord[T any](ctx context.Context, recordURL, apiKey, description string) (T, error) {
	record, err := openrouter.FetchCatalogRecord[T](ctx, recordURL, apiKey)
	if err != nil {
		return fun.Zero[T](), fmt.Errorf("openrouter: fetch %s capabilities: %w", description, err)
	}
	return record, nil
}

func rawyaml(ctx context.Context, spec models.Spec, apiKey string) error {
	endpoints, err := fetchRecord[any](
		ctx, modelURL(spec, "models")+"/endpoints", apiKey, "models endpoints",
	)
	if err != nil {
		return err
	}

	model, err := fetchRecord[any](ctx, modelURL(spec, "model"), apiKey, "model")
	if err != nil {
		return err
	}

	prettyPrint(map[string]any{
		"model":     model,
		"endpoints": endpoints,
	})
	return nil
}

func overview(ctx context.Context, spec models.Spec, apiKey string) error {
	modelResponse, err := fetchRecord[openrouter.GetModelResponse](
		ctx, modelURL(spec, "model"), apiKey, "model",
	)
	if err != nil {
		return err
	}

	model := modelResponse.Data
	prettyPrint(struct {
		CanonicalSlug    any `yaml:"canonical_slug"`
		Name             any `yaml:"name"`
		Pricing          any `yaml:"pricing"`
		ContextLength    any `yaml:"context_length"`
		InputModalities  any `yaml:"input_modalities"`
		OutputModalities any `yaml:"output_modalities"`
		Description      any `yaml:"description"`
		AliasTarget      any `yaml:"alias_target"`
		Reasoning        any `yaml:"reasoning"`
	}{
		CanonicalSlug:    model.CanonicalSlug,
		Name:             model.Name,
		Pricing:          model.Pricing,
		ContextLength:    model.ContextLength,
		InputModalities:  model.Architecture.InputModalities,
		OutputModalities: model.Architecture.OutputModalities,
		Description:      model.Description,
		AliasTarget:      model.AliasTarget,
		Reasoning:        model.Reasoning,
	})
	return nil
}
