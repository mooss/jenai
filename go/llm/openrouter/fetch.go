package openrouter

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/jenai/go/models"
)

func FetchModel(
	ctx context.Context, spec models.Spec, apiKey string,
) (GetModelResponse, error) {
	recordURL := "https://openrouter.ai/api/v1/model/" +
		url.PathEscape(spec.Author) + "/" + url.PathEscape(spec.Model)
	record, err := FetchCatalogRecord[GetModelResponse](
		ctx, recordURL, apiKey,
	)
	if err != nil {
		return fun.Zero[GetModelResponse](), fmt.Errorf(
			"openrouter: fetch model capabilities: %w", err,
		)
	}
	return record, nil
}

func FetchEndpoints(
	ctx context.Context, spec models.Spec, apiKey string,
) (ListEndpointsResponse, error) {
	recordURL := "https://openrouter.ai/api/v1/models/" +
		url.PathEscape(spec.Author) + "/" + url.PathEscape(spec.Model) + "/endpoints"
	record, err := FetchCatalogRecord[ListEndpointsResponse](
		ctx, recordURL, apiKey,
	)
	if err != nil {
		return fun.Zero[ListEndpointsResponse](), fmt.Errorf(
			"openrouter: fetch models endpoints capabilities: %w", err,
		)
	}
	return record, nil
}
