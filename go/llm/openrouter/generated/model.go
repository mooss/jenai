//revive:disable:var-naming,line-length-limit
package generated

// Generated with arguments: ["../scripts/vibe/oapi2go.py", "--package", "generated", "--keep", "200"]

// Single model response
// Returns the model details
type GetModelResponseX200 struct {
	// Information about an AI model available on OpenRouter
	// Example: {
	//   "architecture": {
	//     "input_modalities": [
	//       "text"
	//     ],
	//     "instruct_type": "chatml",
	//     "modality": "text->text",
	//     "output_modalities": [
	//       "text"
	//     ],
	//     "tokenizer": "GPT"
	//   },
	//   "canonical_slug": "openai/gpt-4",
	//   "context_length": 8192,
	//   "created": 1692901234,
	//   "default_parameters": null,
	//   "description": "GPT-4 is a large multimodal model that can solve difficult problems with greater accuracy.",
	//   "expiration_date": null,
	//   "id": "openai/gpt-4",
	//   "knowledge_cutoff": null,
	//   "links": {
	//     "details": "/api/v1/models/openai/gpt-5.4/endpoints"
	//   },
	//   "name": "GPT-4",
	//   "per_request_limits": null,
	//   "pricing": {
	//     "completion": "0.00006",
	//     "image": "0",
	//     "prompt": "0.00003",
	//     "request": "0"
	//   },
	//   "reasoning": {
	//     "default_effort": "medium",
	//     "default_enabled": true,
	//     "mandatory": false,
	//     "supported_efforts": [
	//       "high",
	//       "medium",
	//       "low",
	//       "minimal"
	//     ]
	//   },
	//   "supported_parameters": [
	//     "temperature",
	//     "top_p",
	//     "max_tokens"
	//   ],
	//   "supported_voices": null,
	//   "top_provider": {
	//     "context_length": 8192,
	//     "is_moderated": true,
	//     "max_completion_tokens": 4096
	//   }
	// }
	Data GetModelResponseX200Data `json:"data" yaml:"data"`
}

// Information about an AI model available on OpenRouter
type GetModelResponseX200Data struct {
	// Concrete model targeted by this tilde-latest alias, when applicable
	// Example: {"name": "Claude Sonnet 4.5", "slug": "anthropic/claude-sonnet-4.5"}
	AliasTarget *GetModelResponseX200DataAliasTarget `json:"alias_target" yaml:"alias_target"`
	// Model architecture information
	// Example: {
	//   "input_modalities": [
	//     "text"
	//   ],
	//   "instruct_type": "chatml",
	//   "modality": "text->text",
	//   "output_modalities": [
	//     "text"
	//   ],
	//   "tokenizer": "GPT"
	// }
	Architecture GetModelResponseX200DataArchitecture `json:"architecture" yaml:"architecture"`
	// Third-party benchmark rankings for this model. Omitted when no benchmark data is available.
	// Example: {
	//   "artificial_analysis": {
	//     "agentic_index": 55.8,
	//     "coding_index": 63.2,
	//     "intelligence_index": 71.4
	//   },
	//   "design_arena": [
	//     {
	//       "arena": "models",
	//       "category": "website",
	//       "elo": 1385.2,
	//       "rank": 5,
	//       "win_rate": 62.5
	//     }
	//   ]
	// }
	Benchmarks *GetModelResponseX200DataBenchmarks `json:"benchmarks" yaml:"benchmarks"`
	// Canonical slug for the model
	// Example: "openai/gpt-4"
	CanonicalSlug string `json:"canonical_slug" yaml:"canonical_slug"`
	// Maximum context length in tokens
	// Example: 8192
	ContextLength *int64 `json:"context_length" yaml:"context_length"`
	// Unix timestamp of when the model was created
	// Example: 1692901234
	Created int64 `json:"created" yaml:"created"`
	// Default parameters for this model
	// Example: {
	//   "frequency_penalty": 0,
	//   "presence_penalty": 0,
	//   "repetition_penalty": 1,
	//   "temperature": 0.7,
	//   "top_k": 0,
	//   "top_p": 0.9
	// }
	DefaultParameters *GetModelResponseX200DataDefaultParameters `json:"default_parameters" yaml:"default_parameters"`
	// Description of the model
	// Example: "GPT-4 is a large multimodal model that can solve difficult problems with greater accuracy."
	Description *string `json:"description" yaml:"description"`
	// The date after which the model may be removed. ISO 8601 date string (YYYY-MM-DD) or null if no expiration.
	// Example: "2025-06-01"
	ExpirationDate *string `json:"expiration_date" yaml:"expiration_date"`
	// Hugging Face model identifier, if applicable
	// Example: "microsoft/DialoGPT-medium"
	HuggingFaceId *string `json:"hugging_face_id" yaml:"hugging_face_id"`
	// Unique identifier for the model
	// Example: "openai/gpt-4"
	Id string `json:"id" yaml:"id"`
	// The date up to which the model was trained on data. ISO 8601 date string (YYYY-MM-DD) or null if unknown.
	// Example: "2024-10-01"
	KnowledgeCutoff *string `json:"knowledge_cutoff" yaml:"knowledge_cutoff"`
	// Related API endpoints and resources for this model.
	// Example: {"details": "/api/v1/models/openai/gpt-5.4/endpoints"}
	Links GetModelResponseX200DataLinks `json:"links" yaml:"links"`
	// Display name of the model
	// Example: "GPT-4"
	Name string `json:"name" yaml:"name"`
	// Per-request token limits
	// Example: {"completion_tokens": 1000, "prompt_tokens": 1000}
	PerRequestLimits *GetModelResponseX200DataPerRequestLimits `json:"per_request_limits" yaml:"per_request_limits"`
	// Pricing information for the model
	// Example: {"completion": "0.00006", "image": "0", "prompt": "0.00003", "request": "0"}
	Pricing GetModelResponseX200DataPricing `json:"pricing" yaml:"pricing"`
	// Reasoning effort configuration. Omitted for non-reasoning models and dynamic router models.
	// Example: {
	//   "default_effort": "medium",
	//   "default_enabled": true,
	//   "mandatory": false,
	//   "supported_efforts": [
	//     "high",
	//     "medium",
	//     "low",
	//     "minimal"
	//   ]
	// }
	Reasoning *GetModelResponseX200DataReasoning `json:"reasoning" yaml:"reasoning"`
	// List of supported parameters for this model
	// Allowed values: [
	//   "temperature",
	//   "top_p",
	//   "top_k",
	//   "min_p",
	//   "top_a",
	//   "frequency_penalty",
	//   "presence_penalty",
	//   "repetition_penalty",
	//   "max_tokens",
	//   "max_completion_tokens",
	//   "logit_bias",
	//   "logprobs",
	//   "top_logprobs",
	//   "prediction",
	//   "seed",
	//   "response_format",
	//   "structured_outputs",
	//   "stop",
	//   "tools",
	//   "tool_choice",
	//   "parallel_tool_calls",
	//   "include_reasoning",
	//   "reasoning",
	//   "reasoning_effort",
	//   "web_search_options",
	//   "verbosity"
	// ]
	// Example: "temperature"
	SupportedParameters []string `json:"supported_parameters" yaml:"supported_parameters"`
	// List of supported voice identifiers for TTS models. Null for non-TTS models.
	SupportedVoices *[]string `json:"supported_voices" yaml:"supported_voices"`
	// Information about the top provider for this model
	// Example: {"context_length": 8192, "is_moderated": true, "max_completion_tokens": 4096}
	TopProvider GetModelResponseX200DataTopProvider `json:"top_provider" yaml:"top_provider"`
}

// Concrete model targeted by this tilde-latest alias, when applicable
type GetModelResponseX200DataAliasTarget struct {
	// Human-readable name of the concrete model targeted by this alias
	Name string `json:"name" yaml:"name"`
	// Routable model ID of the concrete target, matching that model row's id
	Slug string `json:"slug" yaml:"slug"`
}

// Model architecture information
type GetModelResponseX200DataArchitecture struct {
	// Supported input modalities
	// Allowed values: ["text", "image", "file", "audio", "video"]
	// Example: "text"
	InputModalities []string `json:"input_modalities" yaml:"input_modalities"`
	// Instruction format type
	// Allowed values: [
	//   "none",
	//   "airoboros",
	//   "alpaca",
	//   "alpaca-modif",
	//   "chatml",
	//   "claude",
	//   "code-llama",
	//   "gemma",
	//   "llama2",
	//   "llama3",
	//   "mistral",
	//   "nemotron",
	//   "neural",
	//   "openchat",
	//   "phi3",
	//   "rwkv",
	//   "vicuna",
	//   "zephyr",
	//   "deepseek-r1",
	//   "deepseek-v3.1",
	//   "qwq",
	//   "qwen3",
	//   null
	// ]
	// Example: "chatml"
	InstructType *string `json:"instruct_type" yaml:"instruct_type"`
	// Primary modality of the model
	// Example: "text->text"
	Modality *string `json:"modality" yaml:"modality"`
	// Supported output modalities
	// Allowed values: ["text", "image", "embeddings", "audio", "video", "rerank", "decisions", "speech", "transcription"]
	// Example: "text"
	OutputModalities []string `json:"output_modalities" yaml:"output_modalities"`
	// Tokenizer type used by the model
	// Allowed values: [
	//   "Router",
	//   "Media",
	//   "Other",
	//   "GPT",
	//   "Claude",
	//   "Gemini",
	//   "Gemma",
	//   "Grok",
	//   "Cohere",
	//   "Nova",
	//   "Qwen",
	//   "Yi",
	//   "DeepSeek",
	//   "Mistral",
	//   "Llama2",
	//   "Llama3",
	//   "Llama4",
	//   "PaLM",
	//   "RWKV",
	//   "Qwen3"
	// ]
	// Example: "GPT"
	Tokenizer *string `json:"tokenizer" yaml:"tokenizer"`
}

// Third-party benchmark rankings for this model. Omitted when no benchmark data is available.
type GetModelResponseX200DataBenchmarks struct {
	// Artificial Analysis benchmark index scores.
	// Example: {"agentic_index": 55.8, "coding_index": 63.2, "intelligence_index": 71.4}
	ArtificialAnalysis *GetModelResponseX200DataBenchmarksArtificialAnalysis `json:"artificial_analysis" yaml:"artificial_analysis"`
	// Design Arena ELO rankings across arena+category pairs.
	// Example: [{"arena": "models", "category": "website", "elo": 1385.2, "rank": 5, "win_rate": 62.5}]
	DesignArena []GetModelResponseX200DataBenchmarksDesignArenaItem `json:"design_arena" yaml:"design_arena"`
}

// Artificial Analysis benchmark index scores.
type GetModelResponseX200DataBenchmarksArtificialAnalysis struct {
	// Artificial Analysis Agentic Index score
	// Example: 55.8
	// format: double
	AgenticIndex *float64 `json:"agentic_index" yaml:"agentic_index"`
	// Artificial Analysis Coding Index score
	// Example: 63.2
	// format: double
	CodingIndex *float64 `json:"coding_index" yaml:"coding_index"`
	// Artificial Analysis Intelligence Index score
	// Example: 71.4
	// format: double
	IntelligenceIndex *float64 `json:"intelligence_index" yaml:"intelligence_index"`
}

// A single Design Arena benchmark entry for a specific arena+category
type GetModelResponseX200DataBenchmarksDesignArenaItem struct {
	// Arena type (e.g. models, builders, agents)
	// Example: "models"
	Arena string `json:"arena" yaml:"arena"`
	// Category within the arena (e.g. website, gamedev, uicomponent)
	// Example: "website"
	Category string `json:"category" yaml:"category"`
	// ELO rating from head-to-head arena battles
	// Example: 1385.2
	// format: double
	Elo float64 `json:"elo" yaml:"elo"`
	// Rank position within this arena+category among models available on OpenRouter (1 = highest ELO)
	// Example: 5
	Rank int64 `json:"rank" yaml:"rank"`
	// Win rate percentage in arena battles
	// Example: 62.5
	// format: double
	WinRate float64 `json:"win_rate" yaml:"win_rate"`
}

// Default parameters for this model
type GetModelResponseX200DataDefaultParameters struct {
	// format: double
	FrequencyPenalty *float64 `json:"frequency_penalty" yaml:"frequency_penalty"`
	// format: double
	PresencePenalty *float64 `json:"presence_penalty" yaml:"presence_penalty"`
	// format: double
	RepetitionPenalty *float64 `json:"repetition_penalty" yaml:"repetition_penalty"`
	// format: double
	Temperature *float64 `json:"temperature" yaml:"temperature"`
	TopK        *int64   `json:"top_k" yaml:"top_k"`
	// format: double
	TopP *float64 `json:"top_p" yaml:"top_p"`
}

// Related API endpoints and resources for this model.
type GetModelResponseX200DataLinks struct {
	// URL for the model details/endpoints API
	// Example: "/api/v1/models/openai/gpt-5.4/endpoints"
	Details string `json:"details" yaml:"details"`
}

// Per-request token limits
type GetModelResponseX200DataPerRequestLimits struct {
	// Maximum completion tokens per request
	// Example: 1000
	CompletionTokens float64 `json:"completion_tokens" yaml:"completion_tokens"`
	// Maximum prompt tokens per request
	// Example: 1000
	PromptTokens float64 `json:"prompt_tokens" yaml:"prompt_tokens"`
}

// Pricing information for the model
type GetModelResponseX200DataPricing struct {
	// Price in USD per audio input token
	Audio *string `json:"audio" yaml:"audio"`
	// Price in USD per audio output token
	AudioOutput *string `json:"audio_output" yaml:"audio_output"`
	// Price in USD per token for completion (output) generation
	Completion string `json:"completion" yaml:"completion"`
	// Fractional discount applied to this endpoint's pricing; the price is multiplied by (1 - discount) (0 = no discount, 1 = free)
	// format: double
	Discount *float64 `json:"discount" yaml:"discount"`
	// Price in USD per input image
	Image *string `json:"image" yaml:"image"`
	// Price in USD per output image
	ImageOutput *string `json:"image_output" yaml:"image_output"`
	// Price in USD per image token
	ImageToken *string `json:"image_token" yaml:"image_token"`
	// Price in USD per cached audio input token
	InputAudioCache *string `json:"input_audio_cache" yaml:"input_audio_cache"`
	// Price in USD per cached input token (read)
	InputCacheRead *string `json:"input_cache_read" yaml:"input_cache_read"`
	// Price per cache-write token, in USD per token. For providers with multiple cache TTLs (e.g. Anthropic), this is the default (5-minute) cache-write rate.
	InputCacheWrite *string `json:"input_cache_write" yaml:"input_cache_write"`
	// Price per 1-hour cache-write token, in USD per token. Only present for providers that price an extended (1-hour) cache TTL separately, such as Anthropic.
	InputCacheWrite1h *string `json:"input_cache_write_1h" yaml:"input_cache_write_1h"`
	// Price in USD per internal reasoning token
	InternalReasoning *string `json:"internal_reasoning" yaml:"internal_reasoning"`
	// Conditional overrides of the base pricing (e.g. long-context or time-based pricing). An entry applies when all of its condition fields (e.g. min_prompt_tokens, or the utc_start/utc_end time window) match the request; among applicable entries, later entries win per key; price keys absent from an entry inherit the base price. The top-level pricing keys always reflect the price that applies under default conditions.
	// Example: {"completion": "0.00002", "min_prompt_tokens": 200000, "prompt": "0.000005"}
	Overrides *[]GetModelResponseX200DataPricingOverridesItem `json:"overrides" yaml:"overrides"`
	// Price in USD per token for prompt (input) processing
	Prompt string `json:"prompt" yaml:"prompt"`
	// Price in USD per request
	Request *string `json:"request" yaml:"request"`
	// Price in USD per web search
	WebSearch *string `json:"web_search" yaml:"web_search"`
}

// A conditional override of the base pricing. An entry applies only when all of its condition fields (e.g. min_prompt_tokens, or the utc_start/utc_end time window) match the request; among applicable entries, later entries win per price key; price keys absent from an entry inherit the base price.
type GetModelResponseX200DataPricingOverridesItem struct {
	// Overridden price in USD per audio input token
	Audio *string `json:"audio" yaml:"audio"`
	// Overridden price in USD per token for completion (output) generation
	Completion *string `json:"completion" yaml:"completion"`
	// Overridden price in USD per cached audio input token
	InputAudioCache *string `json:"input_audio_cache" yaml:"input_audio_cache"`
	// Overridden price in USD per cached input token (read)
	InputCacheRead *string `json:"input_cache_read" yaml:"input_cache_read"`
	// Overridden price in USD per cache-write token
	InputCacheWrite *string `json:"input_cache_write" yaml:"input_cache_write"`
	// Overridden price in USD per 1-hour cache-write token
	InputCacheWrite1h *string `json:"input_cache_write_1h" yaml:"input_cache_write_1h"`
	// Condition: the entry applies when the total prompt tokens of a request are strictly greater than this threshold
	// format: double
	MinPromptTokens *float64 `json:"min_prompt_tokens" yaml:"min_prompt_tokens"`
	// Overridden price in USD per token for prompt (input) processing
	Prompt *string `json:"prompt" yaml:"prompt"`
	// Condition: UTC weekdays the entry applies on, evaluated at the request instant. Scopes the utc_start/utc_end window (or, without a window, the whole UTC day) to the listed days. Absent means every day.
	// Allowed values: ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
	// minItems: 1
	UtcDays *[]string `json:"utc_days" yaml:"utc_days"`
	// Condition: exclusive end of a daily UTC time window as an HHMM clock number (e.g. 400 = 04:00)
	// format: double
	UtcEnd *float64 `json:"utc_end" yaml:"utc_end"`
	// Condition: inclusive start of a daily UTC time window as an HHMM clock number (e.g. 100 = 01:00, 1030 = 10:30). The entry applies while the current UTC time is inside the half-open window [utc_start, utc_end), which may wrap past midnight (utc_start > utc_end).
	// format: double
	UtcStart *float64 `json:"utc_start" yaml:"utc_start"`
}

// Reasoning effort configuration. Omitted for non-reasoning models and dynamic router models.
type GetModelResponseX200DataReasoning struct {
	// Default reasoning effort when the client enables reasoning without specifying effort. Maps to `reasoning.effort` in chat requests. When `"none"`, prefer omitting effort unless the user explicitly disables reasoning.
	DefaultEffort *string `json:"default_effort" yaml:"default_effort"`
	// Default reasoning enabled state when the client does not set `reasoning.enabled`.
	DefaultEnabled *bool `json:"default_enabled" yaml:"default_enabled"`
	// When true, reasoning cannot be disabled and effort "none" is rejected.
	Mandatory bool `json:"mandatory" yaml:"mandatory"`
	// Allowed reasoning effort values for this model, in descending effort order (highest first). Null means no allowlist — all gateway effort values are accepted.
	// Allowed values: ["max", "xhigh", "high", "medium", "low", "minimal", "none", null]
	// Example: "medium"
	SupportedEfforts *[]*string `json:"supported_efforts" yaml:"supported_efforts"`
	// Present and `true` when the model accepts `reasoning.max_tokens` in requests (Anthropic-style) instead of or in addition to `reasoning.effort`. Omitted otherwise.
	SupportsMaxTokens *bool `json:"supports_max_tokens" yaml:"supports_max_tokens"`
}

// Information about the top provider for this model
type GetModelResponseX200DataTopProvider struct {
	// Context length from the top provider
	// Example: 8192
	ContextLength *int64 `json:"context_length" yaml:"context_length"`
	// Whether the top provider moderates content
	// Example: true
	IsModerated bool `json:"is_moderated" yaml:"is_moderated"`
	// Maximum completion tokens from the top provider. Input and output tokens share the context window, so the effective maximum output for a request is further limited by the context remaining after input tokens.
	// Example: 4096
	MaxCompletionTokens *int64 `json:"max_completion_tokens" yaml:"max_completion_tokens"`
}
