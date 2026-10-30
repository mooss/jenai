package main

import (
	"math/big"
	"strings"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/jenai/go/llm/openrouter/generated"
)

// pricingOverview shows text token prices in cents per million tokens (¢/1M).
//
//revive:disable:line-length-limit
type pricingOverview struct {
	Prompt            *string `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	Completion        *string `json:"completion,omitempty" yaml:"completion,omitempty"`
	InputCacheRead    *string `json:"input_cache_read,omitempty" yaml:"input_cache_read,omitempty"`
	InputCacheWrite   *string `json:"input_cache_write,omitempty" yaml:"input_cache_write,omitempty"`
	InputCacheWrite1h *string `json:"input_cache_write_1h,omitempty" yaml:"input_cache_write_1h,omitempty"`
	InternalReasoning *string `json:"internal_reasoning,omitempty" yaml:"internal_reasoning,omitempty"`
}

func textPricing(p generated.GetModelResponseX200DataPricing) *pricingOverview {
	deref := fun.ZeroDeref[string]
	return &pricingOverview{
		Prompt:            centsPerMillion(p.Prompt),
		Completion:        centsPerMillion(p.Completion),
		InputCacheRead:    centsPerMillion(deref(p.InputCacheRead)),
		InputCacheWrite:   centsPerMillion(deref(p.InputCacheWrite)),
		InputCacheWrite1h: centsPerMillion(deref(p.InputCacheWrite1h)),
		InternalReasoning: centsPerMillion(deref(p.InternalReasoning)),
	}
}

// centsPerMillion converts a price in USD per token into a ¢/1M tokens string, or returns nil for
// missing/invalid prices.
func centsPerMillion(price string) *string {
	rat, ok := new(big.Rat).SetString(price)
	if !ok {
		return nil
	}
	rat.Mul(rat, big.NewRat(100_000_000, 1)) // USD/token -> ¢/1M tokens
	s := strings.TrimRight(strings.TrimRight(rat.FloatString(4), "0"), ".")
	if s == "" {
		return nil
	}

	s += " ¢/1M"
	return &s
}
