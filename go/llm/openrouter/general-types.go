//revive:disable:line-length-limit
package openrouter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ModelId contains the required parameters to identify a model.
type ModelId struct {
	Author string `json:"author" yaml:"author"`
	Slug   string `json:"slug" yaml:"slug"`
}

// OpenRouterHeaders contains the optional app-identification headers for an operation.
type OpenRouterHeaders struct {
	HTTPReferer          string `json:"HTTP-Referer,omitempty" yaml:"HTTP-Referer,omitempty"`
	OpenRouterTitle      string `json:"X-OpenRouter-Title,omitempty" yaml:"X-OpenRouter-Title,omitempty"`
	OpenRouterCategories string `json:"X-OpenRouter-Categories,omitempty" yaml:"X-OpenRouter-Categories,omitempty"`
}

// APIErrorResponse is used for the 403, 404, and 500 response bodies.
type APIErrorResponse struct {
	Error              APIError       `json:"error" yaml:"error"`
	OpenRouterMetadata map[string]any `json:"openrouter_metadata,omitempty" yaml:"openrouter_metadata,omitempty"`
	UserID             *string        `json:"user_id,omitempty" yaml:"user_id,omitempty"`
}

type APIError struct {
	Code     int            `json:"code" yaml:"code"`
	Message  string         `json:"message" yaml:"message"`
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxCatalogResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read catalog error response: %w", err)
	}

	apiErr := APIError{Code: resp.StatusCode, Message: resp.Status}
	var response APIErrorResponse
	if json.Unmarshal(data, &response) == nil && response.Error.Message != "" {
		apiErr = response.Error
	}

	return fmt.Errorf("OpenRouter API error (%d): %s", apiErr.Code, apiErr.Message)
}
