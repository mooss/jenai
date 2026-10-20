package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mooss/bagend/go/codec"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
)

func runCapa(cfg *config.Jenai, args []string, source string) error {
	spec, err := capaModelSpec(args, source)
	if err != nil {
		return err
	}

	capabilities, err := openrouter.FetchCapabilities(
		context.Background(), spec, cfg.APIKeys[spec.Provider],
	)
	if err != nil {
		return err
	}

	model, endpoints, err := decodeCapabilities(capabilities, source == "endpoints")
	if err != nil {
		return err
	}

	if source == "endpoints" {
		fmt.Println(pretty(endpoints))
		return nil
	}

	fmt.Println(pretty(map[string]any{
		"model":     model,
		"endpoints": endpoints,
	}))
	return nil
}

func capaModelSpec(args []string, source string) (models.Spec, error) {
	if len(args) != 1 {
		return models.Spec{}, errors.New("usage: jenai capa MODEL")
	}
	if source != "overview" && source != "endpoints" {
		return models.Spec{}, fmt.Errorf("unsupported source %q (expected overview or endpoints)", source)
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

func decodeCapabilities(
	capabilities openrouter.Capabilities, endpointOnly bool,
) (any, any, error) {
	modelResponse, err := codec.FromJSONOpt[
		openrouter.GetModelResponse, codec.WarnUnknownFields,
	](capabilities.Model)
	if err != nil {
		return nil, nil, err
	}

	if endpointOnly {
		endpointsResponse, err := codec.FromJSONOpt[
			openrouter.ListEndpointsResponse, codec.WarnUnknownFields,
		](capabilities.Endpoints)
		if err != nil {
			return nil, nil, err
		}
		return modelResponse.Data, endpointsResponse.Data.Endpoints, nil
	}

	endpointsResponse, err := codec.FromJSON[any](capabilities.Endpoints)
	if err != nil {
		return nil, nil, err
	}

	endpoints, err := filterEndpointCapabilities(endpointsResponse, false)
	if err != nil {
		return nil, nil, err
	}

	return modelResponse.Data, endpoints, nil
}

func filterEndpointCapabilities(response any, onlyEndpoints bool) (any, error) {
	responseMap, ok := response.(map[string]any)
	if !ok {
		return nil, errors.New("invalid endpoint response: expected an object")
	}
	data, ok := responseMap["data"].(map[string]any)
	if !ok {
		return nil, errors.New("invalid endpoint response: missing data object")
	}
	endpoints, ok := data["endpoints"]
	if !ok {
		return nil, errors.New("invalid endpoint response: missing data.endpoints")
	}
	if onlyEndpoints {
		return endpoints, nil
	}

	delete(data, "endpoints")
	return responseMap, nil
}
