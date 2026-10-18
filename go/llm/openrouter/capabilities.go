package openrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mooss/jenai/go/models"
)

// Capabilities contains the model and endpoint records from the OpenRouter catalog.
type Capabilities struct {
	Model     json.RawMessage `json:"model"`
	Endpoints json.RawMessage `json:"endpoints"`
}

// FetchCapabilities retrieves the model-level and endpoint-level catalog records.
// An empty API key makes anonymous catalog requests.
func FetchCapabilities(
	ctx context.Context, spec models.Spec, apiKey string,
) (Capabilities, error) {
	model, err := fetchCatalogRecord(ctx, catalogURL("model", spec), apiKey)
	if err != nil {
		return Capabilities{}, fmt.Errorf("openrouter: fetch model capabilities: %w", err)
	}

	endpointsURL := catalogURL("models", spec) + "/endpoints"
	endpoints, err := fetchCatalogRecord(ctx, endpointsURL, apiKey)
	if err != nil {
		return Capabilities{}, fmt.Errorf("openrouter: fetch endpoint capabilities: %w", err)
	}

	return Capabilities{Model: model, Endpoints: endpoints}, nil
}

func catalogURL(resource string, spec models.Spec) string {
	author := url.PathEscape(spec.Author)
	slug := url.PathEscape(spec.Model)
	return "https://openrouter.ai/api/v1/" + resource + "/" + author + "/" + slug
}

const maxCatalogResponseBytes int64 = 10 << 20 // 10 MiB.

func fetchCatalogRecord(
	ctx context.Context, endpoint, apiKey string,
) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if strings.TrimSpace(apiKey) != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxCatalogResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read catalog response: %w", err)
	}
	if int64(len(data)) > maxCatalogResponseBytes {
		return nil, fmt.Errorf("catalog response exceeds %d bytes", maxCatalogResponseBytes)
	}
	if !json.Valid(data) {
		return nil, errors.New("invalid JSON catalog response")
	}

	return json.RawMessage(data), nil
}
