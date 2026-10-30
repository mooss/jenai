package openrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/jenai/go/models"
)

//////////////
// Typedefs //

const chatCompletionsURL = "https://openrouter.ai/api/v1/chat/completions"

// maxInvalidChunkPreviewLen caps the payload preview included in parse errors to avoid noisy logs
// and leaking large response content.
const maxInvalidChunkPreviewLen = 200

// Chunk is a single SSE `ChatStreamChunk` received from the OpenRouter Chat Completions API.
type Chunk struct {
	// Choices contains the streamed response choices in this chunk.
	Choices []Choice `json:"choices"`
	// Usage contains token usage information when it is included in the chunk.
	Usage *Usage `json:"usage,omitempty"`
	// Error contains an API error returned in the stream, if any.
	Error *apiError `json:"error,omitempty"`
}

// Choice holds one streamed delta and its completion status.
type Choice struct {
	// Delta contains the incremental response content for this choice.
	Delta Delta `json:"delta"`
	// FinishReason is nil while the choice is still streaming and non-nil when it finishes.
	FinishReason *string `json:"finish_reason"`
}

// Delta carries either answer text or reasoning text.
// They must stay separate and never be concatenated.
type Delta struct {
	// Role identifies the speaker, typically "assistant" in the initial chunk.
	Role string `json:"role,omitempty"`
	// Content is an incremental piece of the answer text.
	Content string `json:"content,omitempty"`
	// Reasoning is an incremental piece of the model's reasoning text.
	Reasoning string `json:"reasoning,omitempty"`
}

// Usage contains token counts and cost information reported by the API.
type Usage struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens,omitempty"`
	// PromptTokensDetails contains details about prompt token usage.
	PromptTokensDetails *PromptTokensDetails `json:"prompt_tokens_details,omitempty"`
	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens,omitempty"`
	// CompletionTokensDetails contains details about completion token usage.
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
	// TotalTokens is the total number of prompt and completion tokens.
	TotalTokens int `json:"total_tokens,omitempty"`
	// Cost is the total cost reported by the API.
	Cost *float64 `json:"cost,omitempty"`
	// CostDetails contains details about the reported cost.
	CostDetails *CostDetails `json:"cost_details,omitempty"`
}

// PromptTokensDetails contains details about prompt token usage.
type PromptTokensDetails struct {
	CachedTokens     int `json:"cached_tokens,omitempty"`
	CacheWriteTokens int `json:"cache_write_tokens,omitempty"`
	AudioTokens      int `json:"audio_tokens,omitempty"`
}

// CompletionTokensDetails contains details about completion token usage.
type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// CostDetails contains details about the reported cost.
type CostDetails struct {
	UpstreamInferenceCost float64 `json:"upstream_inference_cost,omitempty"`
}

// apiError represents an error returned in a streamed API chunk.
type apiError struct {
	// Code is the API error code, when provided.
	Code int `json:"code,omitempty"`
	// Message describes the API error.
	Message string `json:"message,omitempty"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("openrouter: code %d: %s", e.Code, e.Message)
}

// chatRequest is the request body sent to the OpenRouter Chat Completions API.
type chatRequest struct {
	// Model is the OpenRouter model identifier.
	Model string `json:"model"`
	// Messages contains the conversation messages sent to the model.
	Messages []chatMessage `json:"messages"`
	// Stream requests that the response be returned as a stream of SSE chunks.
	Stream bool `json:"stream"`
	// Reasoning configures the model's reasoning effort and summary verbosity.
	Reasoning chatReasoning `json:"reasoning"`
}

// chatMessage is a message in a Chat Completions request.
type chatMessage struct {
	// Role identifies the speaker, such as "user" or "assistant".
	Role string `json:"role"`
	// Content contains the message text.
	Content string `json:"content"`
}

// chatReasoning configures reasoning generation for a Chat Completions request.
type chatReasoning struct {
	// Effort sets the reasoning effort, such as "medium".
	Effort string `json:"effort,omitempty"`
	// Summary sets the reasoning summary verbosity, such as "auto".
	Summary string `json:"summary"`
}

////////////
// Stream //

// StreamCompletion starts a streaming completion and returns the raw chunk and error channels.
func StreamCompletion(
	ctx context.Context, spec models.Spec, apiKey, prompt string,
) (<-chan Chunk, <-chan error, error) {
	req, err := newRequest(ctx, spec, apiKey, prompt)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if err := checkStatus(resp); err != nil {
		resp.Body.Close()
		return nil, nil, err
	}

	chunks, errs := scan(ctx, resp)
	return chunks, errs, nil
}

func newRequest(
	ctx context.Context, spec models.Spec, apiKey, prompt string,
) (*http.Request, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("openrouter: missing API key")
	}
	if spec.Provider != "openrouter" {
		msg := "openrouter: unsupported provider %q: expected %q"
		return nil, fmt.Errorf(msg, spec.Provider, "openrouter")
	}

	effort := fun.ZeroDeref(spec.Reasoning).Effort

	body, err := json.Marshal(chatRequest{
		Model:     spec.ID(),
		Messages:  []chatMessage{{Role: "user", Content: prompt}},
		Stream:    true,
		Reasoning: chatReasoning{Effort: effort, Summary: "auto"},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chatCompletionsURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	return req, nil
}

//////////////
// SSE Scan //

func scan(ctx context.Context, resp *http.Response) (<-chan Chunk, <-chan error) {
	chunks := make(chan Chunk)
	errs := make(chan error, 1)

	go func() {
		defer close(chunks)
		defer close(errs)
		defer resp.Body.Close()

		scanStream(ctx, resp.Body, chunks, errs)
	}()

	return chunks, errs
}

func scanStream(ctx context.Context, body io.Reader, chunks chan<- Chunk, errs chan<- error) {
	reader := bufio.NewReader(body)

	for {
		line, readErr := reader.ReadString('\n')
		if !processSSELine(ctx, line, chunks, errs) {
			return
		}
		if readErr != nil {
			reportStreamReadError(ctx, readErr, errs)
			return
		}
	}
}

func processSSELine(ctx context.Context, line string, chunks chan<- Chunk, errs chan<- error) bool {
	payload, ok := ssePayload(line)
	if !ok {
		return true
	}
	if payload == "[DONE]" {
		return false
	}

	chunk, err := parseChunk(payload)
	if err != nil {
		errs <- err
		return false
	}
	if chunk.Error != nil {
		errs <- chunk.Error
		return false
	}

	select {
	case chunks <- chunk:
		return true
	case <-ctx.Done():
		errs <- ctx.Err()
		return false
	}
}

func parseChunk(payload string) (Chunk, error) {
	var chunk Chunk
	if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
		preview := truncatePreview(payload, maxInvalidChunkPreviewLen)
		return Chunk{}, fmt.Errorf("openrouter: invalid chunk %q: %w", preview, err)
	}
	return chunk, nil
}

func reportStreamReadError(ctx context.Context, readErr error, errs chan<- error) {
	if ctx.Err() != nil {
		errs <- ctx.Err()
		return
	}
	if !errors.Is(readErr, io.EOF) {
		errs <- fmt.Errorf("openrouter: stream read: %w", readErr)
	}
}

/////////////////
// SSE Parsing //

func ssePayload(line string) (string, bool) {
	if !strings.HasPrefix(line, "data:") {
		return "", false
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "" {
		return "", false
	}
	return payload, true
}

func truncatePreview(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	// Truncate on rune boundary to avoid splitting UTF-8.
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "… [truncated]"
}
