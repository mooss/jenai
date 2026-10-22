package openrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mooss/bagend/go/codec"
)

func FetchCatalogRecord[T any](
	ctx context.Context, endpoint, apiKey string,
) (T, error) {
	data, err := fetchCatalogRaw(ctx, endpoint, apiKey)
	if err != nil {
		var zero T
		return zero, err
	}

	return codec.FromJSONOpt[T, codec.WarnUnknownFields](data)
}

const maxResponseBytes int64 = 10 << 20 // 10 MiB.

func fetchCatalogRaw(ctx context.Context, endpoint, apiKey string) ([]byte, error) {
	req, err := newCatalogRequest(ctx, endpoint, apiKey)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp); err != nil {
		return nil, err
	}

	return readCatalogResponse(resp.Body)
}

func newCatalogRequest(ctx context.Context, endpoint, apiKey string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	setCatalogRequestHeaders(req, apiKey)

	return req, nil
}

func setCatalogRequestHeaders(req *http.Request, apiKey string) {
	headers := Headers{}
	if headers.HTTPReferer != "" {
		req.Header.Set("HTTP-Referer", headers.HTTPReferer)
	}
	if headers.OpenRouterTitle != "" {
		req.Header.Set("X-OpenRouter-Title", headers.OpenRouterTitle)
	}
	if headers.OpenRouterCategories != "" {
		req.Header.Set("X-OpenRouter-Categories", headers.OpenRouterCategories)
	}
	if strings.TrimSpace(apiKey) != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
}

func readCatalogResponse(body io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(body, maxResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read catalog response: %w", err)
	}
	if int64(len(data)) > maxResponseBytes {
		return nil, fmt.Errorf("catalog response exceeds %d bytes", maxResponseBytes)
	}
	if !json.Valid(data) {
		return nil, errors.New("invalid JSON catalog response")
	}

	return data, nil
}
