package main

import (
	"context"
	"errors"
	"fmt"

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
	case "endpoints":
		endpointsResponse, err := openrouter.FetchEndpoints(ctx, spec, apiKey)
		if err != nil {
			return err
		}
		prettyPrint(endpointsResponse.Data.Endpoints)
		return nil
	}

	modelResponse, err := openrouter.FetchModel(ctx, spec, apiKey)
	if err != nil {
		return err
	}

	prettyPrint(modelResponse.Data)
	return nil
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

func rawyaml(ctx context.Context, spec models.Spec, apiKey string) error {
	model, err := openrouter.FetchModel(ctx, spec, apiKey)
	if err != nil {
		return err
	}

	endpoints, err := openrouter.FetchEndpoints(ctx, spec, apiKey)
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
	modelResponse, err := openrouter.FetchModel(ctx, spec, apiKey)
	if err != nil {
		return err
	}

	model := modelResponse.Data
	//revive:disable:line-length-limit
	prettyPrint(struct {
		CanonicalSlug    any                         `json:"canonical_slug,omitempty" yaml:"canonical_slug,omitempty"`
		Name             any                         `json:"name,omitempty" yaml:"name,omitempty"`
		Pricing          *openrouter.PricingOverview `json:"pricing,omitempty" yaml:"pricing,omitempty"`
		ContextLength    any                         `json:"context_length,omitempty" yaml:"context_length,omitempty"`
		InputModalities  any                         `json:"input_modalities,omitempty" yaml:"input_modalities,omitempty"`
		OutputModalities any                         `json:"output_modalities,omitempty" yaml:"output_modalities,omitempty"`
		Description      any                         `json:"description,omitempty" yaml:"description,omitempty"`
		AliasTarget      any                         `json:"alias_target,omitempty" yaml:"alias_target,omitempty"`
		Reasoning        any                         `json:"reasoning,omitempty" yaml:"reasoning,omitempty"`
	}{
		CanonicalSlug:    model.CanonicalSlug,
		Name:             model.Name,
		Pricing:          openrouter.TextPricing(model.Pricing),
		ContextLength:    model.ContextLength,
		InputModalities:  model.Architecture.InputModalities,
		OutputModalities: model.Architecture.OutputModalities,
		Description:      model.Description,
		AliasTarget:      model.AliasTarget,
		Reasoning:        model.Reasoning,
	})
	return nil
}
