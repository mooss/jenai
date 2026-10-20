//revive:disable:line-length-limit
package openrouter

// GetModelResponse is the successful response body.
type GetModelResponse struct {
	Data Model `json:"data" yaml:"data"`
}

type Model struct {
	ID                  string             `json:"id" yaml:"id"`
	CanonicalSlug       string             `json:"canonical_slug" yaml:"canonical_slug"`
	Name                string             `json:"name" yaml:"name"`
	Created             int64              `json:"created" yaml:"created"`
	Pricing             ModelPricing       `json:"pricing" yaml:"pricing"`
	ContextLength       *int64             `json:"context_length" yaml:"context_length"`
	Architecture        ModelArchitecture  `json:"architecture" yaml:"architecture"`
	TopProvider         ModelTopProvider   `json:"top_provider" yaml:"top_provider"`
	PerRequestLimits    *PerRequestLimits  `json:"per_request_limits" yaml:"per_request_limits"`
	SupportedParameters []string           `json:"supported_parameters" yaml:"supported_parameters"`
	DefaultParameters   *DefaultParameters `json:"default_parameters" yaml:"default_parameters"`
	SupportedVoices     []string           `json:"supported_voices" yaml:"supported_voices"`
	Links               ModelLinks         `json:"links" yaml:"links"`

	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	ExpirationDate  *string           `json:"expiration_date,omitempty" yaml:"expiration_date,omitempty"`
	KnowledgeCutoff *string           `json:"knowledge_cutoff,omitempty" yaml:"knowledge_cutoff,omitempty"`
	HuggingFaceID   *string           `json:"hugging_face_id,omitempty" yaml:"hugging_face_id,omitempty"`
	AliasTarget     *ModelAliasTarget `json:"alias_target,omitempty" yaml:"alias_target,omitempty"`
	Benchmarks      *ModelBenchmarks  `json:"benchmarks,omitempty" yaml:"benchmarks,omitempty"`
	Reasoning       *ModelReasoning   `json:"reasoning,omitempty" yaml:"reasoning,omitempty"`
}

type ModelAliasTarget struct {
	Slug string `json:"slug" yaml:"slug"`
	Name string `json:"name" yaml:"name"`
}

type ModelArchitecture struct {
	Modality         *string  `json:"modality" yaml:"modality"`
	InputModalities  []string `json:"input_modalities" yaml:"input_modalities"`
	OutputModalities []string `json:"output_modalities" yaml:"output_modalities"`
	Tokenizer        string   `json:"tokenizer,omitempty" yaml:"tokenizer,omitempty"`
	InstructType     *string  `json:"instruct_type,omitempty" yaml:"instruct_type,omitempty"`
}

type ModelTopProvider struct {
	IsModerated         bool   `json:"is_moderated" yaml:"is_moderated"`
	ContextLength       *int64 `json:"context_length" yaml:"context_length"`
	MaxCompletionTokens *int64 `json:"max_completion_tokens" yaml:"max_completion_tokens"`
}

type PerRequestLimits struct {
	PromptTokens     float64 `json:"prompt_tokens" yaml:"prompt_tokens"`
	CompletionTokens float64 `json:"completion_tokens" yaml:"completion_tokens"`
}

type DefaultParameters struct {
	FrequencyPenalty  *float64 `json:"frequency_penalty,omitempty" yaml:"frequency_penalty,omitempty"`
	PresencePenalty   *float64 `json:"presence_penalty,omitempty" yaml:"presence_penalty,omitempty"`
	RepetitionPenalty *float64 `json:"repetition_penalty,omitempty" yaml:"repetition_penalty,omitempty"`
	Temperature       *float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	TopK              *int64   `json:"top_k,omitempty" yaml:"top_k,omitempty"`
	TopP              *float64 `json:"top_p,omitempty" yaml:"top_p,omitempty"`
}

type ModelLinks struct {
	Details string `json:"details" yaml:"details"`
}

type ModelPricing struct {
	Prompt            string            `json:"prompt" yaml:"prompt"`
	Completion        string            `json:"completion" yaml:"completion"`
	Request           string            `json:"request,omitempty" yaml:"request,omitempty"`
	Image             string            `json:"image,omitempty" yaml:"image,omitempty"`
	Audio             string            `json:"audio,omitempty" yaml:"audio,omitempty"`
	AudioOutput       string            `json:"audio_output,omitempty" yaml:"audio_output,omitempty"`
	ImageOutput       string            `json:"image_output,omitempty" yaml:"image_output,omitempty"`
	ImageToken        string            `json:"image_token,omitempty" yaml:"image_token,omitempty"`
	InputAudioCache   string            `json:"input_audio_cache,omitempty" yaml:"input_audio_cache,omitempty"`
	InputCacheRead    string            `json:"input_cache_read,omitempty" yaml:"input_cache_read,omitempty"`
	InputCacheWrite   string            `json:"input_cache_write,omitempty" yaml:"input_cache_write,omitempty"`
	InputCacheWrite1H string            `json:"input_cache_write_1h,omitempty" yaml:"input_cache_write_1h,omitempty"`
	InternalReasoning string            `json:"internal_reasoning,omitempty" yaml:"internal_reasoning,omitempty"`
	WebSearch         string            `json:"web_search,omitempty" yaml:"web_search,omitempty"`
	Discount          *float64          `json:"discount,omitempty" yaml:"discount,omitempty"`
	Overrides         []PricingOverride `json:"overrides,omitempty" yaml:"overrides,omitempty"`
}

type PricingOverride struct {
	Prompt            string   `json:"prompt,omitempty" yaml:"prompt,omitempty"`
	Completion        string   `json:"completion,omitempty" yaml:"completion,omitempty"`
	Audio             string   `json:"audio,omitempty" yaml:"audio,omitempty"`
	InputAudioCache   string   `json:"input_audio_cache,omitempty" yaml:"input_audio_cache,omitempty"`
	InputCacheRead    string   `json:"input_cache_read,omitempty" yaml:"input_cache_read,omitempty"`
	InputCacheWrite   string   `json:"input_cache_write,omitempty" yaml:"input_cache_write,omitempty"`
	InputCacheWrite1H string   `json:"input_cache_write_1h,omitempty" yaml:"input_cache_write_1h,omitempty"`
	MinPromptTokens   *float64 `json:"min_prompt_tokens,omitempty" yaml:"min_prompt_tokens,omitempty"`
	UTCDays           []string `json:"utc_days,omitempty" yaml:"utc_days,omitempty"`
	UTCStart          *float64 `json:"utc_start,omitempty" yaml:"utc_start,omitempty"`
	UTCEnd            *float64 `json:"utc_end,omitempty" yaml:"utc_end,omitempty"`
}

type ModelReasoning struct {
	Mandatory         bool      `json:"mandatory" yaml:"mandatory"`
	DefaultEffort     *string   `json:"default_effort,omitempty" yaml:"default_effort,omitempty"`
	DefaultEnabled    *bool     `json:"default_enabled,omitempty" yaml:"default_enabled,omitempty"`
	SupportedEfforts  []*string `json:"supported_efforts,omitempty" yaml:"supported_efforts,omitempty"`
	SupportsMaxTokens *bool     `json:"supports_max_tokens,omitempty" yaml:"supports_max_tokens,omitempty"`
}

type ModelBenchmarks struct {
	ArtificialAnalysis *ArtificialAnalysisBenchmarks `json:"artificial_analysis,omitempty" yaml:"artificial_analysis,omitempty"`
	DesignArena        []DesignArenaBenchmark        `json:"design_arena" yaml:"design_arena"`
}

type ArtificialAnalysisBenchmarks struct {
	IntelligenceIndex *float64 `json:"intelligence_index" yaml:"intelligence_index"`
	CodingIndex       *float64 `json:"coding_index" yaml:"coding_index"`
	AgenticIndex      *float64 `json:"agentic_index" yaml:"agentic_index"`
}

type DesignArenaBenchmark struct {
	Arena    string  `json:"arena" yaml:"arena"`
	Category string  `json:"category" yaml:"category"`
	ELO      float64 `json:"elo" yaml:"elo"`
	WinRate  float64 `json:"win_rate" yaml:"win_rate"`
	Rank     int     `json:"rank" yaml:"rank"`
}
