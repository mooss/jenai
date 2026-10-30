package openrouter

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mooss/bagend/go/obj/osl"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/renderer"
)

func NewStreamer(
	ctx context.Context, spec models.Spec,
	apiKey, prompt string,
	reader renderer.CompletionReader,
) (*Streamer, error) {
	chunks, errs, err := CreateCompletionStream(ctx, spec, apiKey, prompt)
	if err != nil {
		return nil, err
	}
	return &Streamer{
		usage:    nil,
		consumer: reader,
		chunks:   chunks,
		errs:     errs,
	}, nil
}

type Streamer struct {
	usage    *Usage
	consumer renderer.CompletionReader
	chunks   <-chan Chunk
	errs     <-chan error
}

func (r *Streamer) ProcessChunks() error {
	for chunk := range r.chunks {
		if err := r.handleChunk(chunk); err != nil {
			return err
		}
	}

	return <-r.errs
}

func (r *Streamer) handleChunk(chunk Chunk) error {
	if chunk.Usage != nil {
		r.usage = chunk.Usage
	}
	if len(chunk.Choices) == 0 {
		if chunk.Usage == nil { // Expected final usage-only chunk.
			fmt.Fprintln(os.Stderr, "warning: chunk with empty choices and no usage")
		}

		return nil
	}

	delta := chunk.Choices[0].Delta
	if delta.Reasoning != "" {
		if err := r.consumer.ReadReasoning(delta.Reasoning); err != nil {
			return err
		}
	}
	if delta.Content != "" {
		if err := r.consumer.ReadContent(delta.Content); err != nil {
			return err
		}
	}
	return nil
}

func (r *Streamer) FormatUsage() string {
	if r.usage == nil {
		return ""
	}
	u := r.usage
	costStr := "¢? (unavailable)"
	if u.Cost != nil {
		costStr = fmt.Sprintf("¢%.4f", 100*(*u.Cost))
	}
	parts := osl.New(
		costStr,
		fmt.Sprintf(":completion %d+%d=%d", u.PromptTokens, u.CompletionTokens, u.TotalTokens),
	)
	if u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0 {
		parts.Push(fmt.Sprintf(":cached %d", u.PromptTokensDetails.CachedTokens))
	}
	if u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0 {
		parts.Push(fmt.Sprintf(":reasoning %d", u.CompletionTokensDetails.ReasoningTokens))
	}

	return fmt.Sprintf("Cost: %s\n", strings.Join(parts, " "))
}
