// Generated with arguments: ["../scripts/vibe/oapi2go.py", "--package", "generated", "--keep", "200"]
//
//revive:disable:var-naming,line-length-limit
package generated

// Returns a list of endpoints
type ListEndpointsResponseX200 struct {
	// List of available endpoints for a model
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
	//   "created": 1692901234,
	//   "description": "GPT-4 is a large multimodal model that can solve difficult problems with greater accuracy.",
	//   "endpoints": [
	//     {
	//       "context_length": 8192,
	//       "latency_last_30m": {
	//         "p50": 0.25,
	//         "p75": 0.35,
	//         "p90": 0.48,
	//         "p99": 0.85
	//       },
	//       "max_completion_tokens": 4096,
	//       "max_prompt_tokens": 8192,
	//       "model_id": "openai/gpt-4",
	//       "model_name": "GPT-4",
	//       "name": "OpenAI: GPT-4",
	//       "pricing": {
	//         "completion": "0.00006",
	//         "image": "0",
	//         "prompt": "0.00003",
	//         "request": "0"
	//       },
	//       "provider_name": "OpenAI",
	//       "quantization": "fp16",
	//       "status": 0,
	//       "supported_parameters": [
	//         "temperature",
	//         "top_p",
	//         "max_tokens",
	//         "frequency_penalty",
	//         "presence_penalty"
	//       ],
	//       "supports_diarization": false,
	//       "supports_implicit_caching": true,
	//       "tag": "openai",
	//       "throughput_last_30m": {
	//         "p50": 45.2,
	//         "p75": 38.5,
	//         "p90": 28.3,
	//         "p99": 15.1
	//       },
	//       "uptime_last_1d": 99.8,
	//       "uptime_last_30m": 99.5,
	//       "uptime_last_5m": 100
	//     }
	//   ],
	//   "id": "openai/gpt-4",
	//   "name": "GPT-4"
	// }
	Data ListEndpointsResponseX200Data `json:"data" yaml:"data"`
}

// List of available endpoints for a model
type ListEndpointsResponseX200Data struct {
	// Model architecture information
	Architecture ListEndpointsResponseX200DataArchitecture `json:"architecture" yaml:"architecture"`
	// Unix timestamp of when the model was created
	// Example: 1692901234
	Created int64 `json:"created" yaml:"created"`
	// Description of the model
	// Example: "GPT-4 is a large multimodal model that can solve difficult problems with greater accuracy."
	Description string `json:"description" yaml:"description"`
	// List of available endpoints for this model
	// Example: {
	//   "context_length": 8192,
	//   "latency_last_30m": {
	//     "p50": 0.25,
	//     "p75": 0.35,
	//     "p90": 0.48,
	//     "p99": 0.85
	//   },
	//   "max_completion_tokens": 4096,
	//   "max_prompt_tokens": 8192,
	//   "model_id": "openai/gpt-4",
	//   "model_name": "GPT-4",
	//   "name": "OpenAI: GPT-4",
	//   "native_tools": {
	//     "openrouter:web_search": {
	//       "type": "web_search_20260209"
	//     }
	//   },
	//   "perf_last_30m_by_workload": {
	//     "text_generation": {
	//       "latency": {
	//         "p50": 250,
	//         "p75": 350,
	//         "p90": 480,
	//         "p99": 850
	//       },
	//       "request_count": 1000,
	//       "throughput": {
	//         "p50": 45.2,
	//         "p75": 38.5,
	//         "p90": 28.3,
	//         "p99": 15.1
	//       }
	//     }
	//   },
	//   "pricing": {
	//     "completion": "0.00006",
	//     "image": "0",
	//     "prompt": "0.00003",
	//     "request": "0"
	//   },
	//   "provider_name": "OpenAI",
	//   "quantization": "fp16",
	//   "service_tier": null,
	//   "status": 0,
	//   "supported_parameters": [
	//     "temperature",
	//     "top_p",
	//     "max_tokens"
	//   ],
	//   "supports_diarization": false,
	//   "supports_image_reference": false,
	//   "supports_implicit_caching": true,
	//   "supports_multiple_audio_references": false,
	//   "supports_tool_choice": {
	//     "auto": true,
	//     "function": true,
	//     "none": true,
	//     "required": true
	//   },
	//   "supports_voice_cloning": false,
	//   "tag": "openai",
	//   "throughput_last_30m": {
	//     "p50": 45.2,
	//     "p75": 38.5,
	//     "p90": 28.3,
	//     "p99": 15.1
	//   },
	//   "uptime_last_1d": 99.8,
	//   "uptime_last_30m": 99.5,
	//   "uptime_last_5m": 100
	// }
	Endpoints []ListEndpointsResponseX200DataEndpointsItem `json:"endpoints" yaml:"endpoints"`
	// Unique identifier for the model
	// Example: "openai/gpt-4"
	Id string `json:"id" yaml:"id"`
	// Display name of the model
	// Example: "GPT-4"
	Name string `json:"name" yaml:"name"`
}

// Model architecture information
type ListEndpointsResponseX200DataArchitecture struct {
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
	// Example: "text"
	Modality *string `json:"modality" yaml:"modality"`
	// Supported output modalities
	// Allowed values: ["text", "image", "embeddings", "audio", "video", "rerank", "decisions", "speech", "transcription"]
	// Example: "text"
	OutputModalities []string `json:"output_modalities" yaml:"output_modalities"`
	// Tokenizer type used by the model
	// anyOf composition:
	// {
	//   "description": "Tokenizer type used by the model",
	//   "enum": [
	//     "Router",
	//     "Media",
	//     "Other",
	//     "GPT",
	//     "Claude",
	//     "Gemini",
	//     "Gemma",
	//     "Grok",
	//     "Cohere",
	//     "Nova",
	//     "Qwen",
	//     "Yi",
	//     "DeepSeek",
	//     "Mistral",
	//     "Llama2",
	//     "Llama3",
	//     "Llama4",
	//     "PaLM",
	//     "RWKV",
	//     "Qwen3"
	//   ],
	//   "example": "GPT",
	//   "type": "string",
	//   "x-speakeasy-unknown-values": "allow"
	// }
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
	Tokenizer string `json:"tokenizer" yaml:"tokenizer"`
}

// Information about a specific model endpoint
type ListEndpointsResponseX200DataEndpointsItem struct {
	ContextLength int64 `json:"context_length" yaml:"context_length"`
	// Latency percentiles in milliseconds over the last 30 minutes. Latency measures time to first token. Only visible when authenticated with an API key or cookie; returns null for unauthenticated requests.
	// Example: {"p50": 25.5, "p75": 35.2, "p90": 48.7, "p99": 85.3}
	LatencyLast30m *ListEndpointsResponseX200DataEndpointsItemLatencyLast30m `json:"latency_last_30m" yaml:"latency_last_30m"`
	// Maximum completion tokens for this endpoint. Input and output tokens share the context window, so the effective maximum output for a request is further limited by the context remaining after input tokens.
	MaxCompletionTokens *int64 `json:"max_completion_tokens" yaml:"max_completion_tokens"`
	MaxPromptTokens     *int64 `json:"max_prompt_tokens" yaml:"max_prompt_tokens"`
	// The unique identifier for the model (permaslug)
	// Example: "openai/gpt-4"
	ModelId   string `json:"model_id" yaml:"model_id"`
	ModelName string `json:"model_name" yaml:"model_name"`
	Name      string `json:"name" yaml:"name"`
	// The server tools this endpoint accepts as the provider's own built-in tool (`engine: "native"`) instead of an OpenRouter engine, keyed by canonical `openrouter:*` name. Each value names the provider tool type the request is translated to. Where that tool runs (provider-side, or returned to the client as with Anthropic bash) is documented per tool. Empty when the provider has none.
	// Example: {"openrouter:web_search": {"type": "web_search_20260209"}}
	NativeTools map[string]ListEndpointsResponseX200DataEndpointsItemNativeToolsValue `json:"native_tools" yaml:"native_tools"`
	// Endpoint performance over the last 30 minutes, keyed by the kind of request served (e.g. `text_generation`, `image_generation`). Additive to the legacy singular latency and throughput fields; image and video generation report end-to-end latency. Only visible when authenticated with an API key or cookie.
	PerfLast30mByWorkload *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkload `json:"perf_last_30m_by_workload" yaml:"perf_last_30m_by_workload"`
	Pricing               ListEndpointsResponseX200DataEndpointsItemPricing                `json:"pricing" yaml:"pricing"`
	// Allowed values: [
	//   "Modal",
	//   "AkashML",
	//   "AI21",
	//   "AionLabs",
	//   "Alibaba",
	//   "Ambient",
	//   "Baidu",
	//   "Amazon Bedrock",
	//   "Amazon Nova",
	//   "Anthropic",
	//   "Arcee AI",
	//   "AssemblyAI",
	//   "AtlasCloud",
	//   "Avian",
	//   "Azure",
	//   "BaseTen",
	//   "BytePlus",
	//   "Black Forest Labs",
	//   "Cerebras",
	//   "Chutes",
	//   "Cirrascale",
	//   "Claude Platform on AWS",
	//   "Clarifai",
	//   "Cloudflare",
	//   "Cohere",
	//   "CoreWeave",
	//   "Cosine",
	//   "Crucible",
	//   "Crusoe",
	//   "Darkbloom",
	//   "Databricks",
	//   "Decart",
	//   "Deepgram",
	//   "DeepInfra",
	//   "DeepSeek",
	//   "DekaLLM",
	//   "DigitalOcean",
	//   "ElevenLabs",
	//   "Featherless",
	//   "Fireworks",
	//   "Fish Audio",
	//   "Friendli",
	//   "GMICloud",
	//   "Google",
	//   "Google AI Studio",
	//   "Groq",
	//   "HeyGen",
	//   "Inception",
	//   "Inceptron",
	//   "InferenceNet",
	//   "Ionstream",
	//   "Infermatic",
	//   "Io Net",
	//   "Inferact vLLM",
	//   "Inflection",
	//   "Liquid",
	//   "Makora",
	//   "Mara",
	//   "Mancer 2",
	//   "Meta",
	//   "Minimax",
	//   "ModelRun",
	//   "Mistral",
	//   "Modular",
	//   "Moonshot AI",
	//   "Morph",
	//   "VoyageAI by MongoDB",
	//   "TypeSafe",
	//   "Near AI",
	//   "Nebius",
	//   "Nex AGI",
	//   "NextBit",
	//   "Novita",
	//   "Nvidia",
	//   "Ollama",
	//   "OpenAI",
	//   "OpenInference",
	//   "Parasail",
	//   "ParetoAlpha",
	//   "Poolside",
	//   "PrimeIntellect",
	//   "Perceptron",
	//   "Perplexity",
	//   "Phala",
	//   "Recraft",
	//   "Recursive",
	//   "Reka",
	//   "Relace",
	//   "Respan",
	//   "Sail Research",
	//   "Sakana AI",
	//   "SambaNova",
	//   "ScaleDown",
	//   "Seed",
	//   "SiliconFlow",
	//   "Sourceful",
	//   "StepFun",
	//   "Stealth",
	//   "StreamLake",
	//   "Switchpoint",
	//   "Tencent",
	//   "Tenstorrent",
	//   "Thinking Machines",
	//   "Together",
	//   "Unbiased",
	//   "Upstage",
	//   "Venice",
	//   "Wafer",
	//   "WandB",
	//   "Quiver",
	//   "Krea",
	//   "Runway",
	//   "Xiaomi",
	//   "xAI",
	//   "Z.AI",
	//   "FakeProvider"
	// ]
	// Example: "OpenAI"
	ProviderName string `json:"provider_name" yaml:"provider_name"`
	// anyOf composition:
	// {
	//   "enum": [
	//     "int4",
	//     "int8",
	//     "fp4",
	//     "mxfp4",
	//     "nvfp4",
	//     "fp6",
	//     "fp8",
	//     "mxfp8",
	//     "fp16",
	//     "bf16",
	//     "fp32",
	//     "unknown"
	//   ],
	//   "example": "fp16",
	//   "type": "string",
	//   "x-speakeasy-unknown-values": "allow"
	// }
	Quantization string `json:"quantization" yaml:"quantization"`
	// The service tier this endpoint sells, or null for a standard endpoint. The same value whether the provider selects the tier with a request parameter or with a dedicated model id. Always the canonical tier: priority covers the fast request alias and the turbo, highspeed, and prime provider model-id spellings, none of which this field reports.
	// Allowed values: ["flex", "priority", "ultrafast", null]
	// Example: "priority"
	ServiceTier *string `json:"service_tier" yaml:"service_tier"`
	// Allowed values: [0, -1, -2, -3, -5, -10]
	// Example: 0
	Status *int64 `json:"status" yaml:"status"`
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
	// Whether this STT endpoint can label each transcript segment with a speaker (speaker diarization). False when the endpoint has not declared support.
	SupportsDiarization bool `json:"supports_diarization" yaml:"supports_diarization"`
	// Whether this TTS endpoint accepts an `image_url` reference describing the desired voice. Requests carrying an image reference are only routed to endpoints where this is true.
	SupportsImageReference  *bool `json:"supports_image_reference" yaml:"supports_image_reference"`
	SupportsImplicitCaching bool  `json:"supports_implicit_caching" yaml:"supports_implicit_caching"`
	// Whether this TTS endpoint accepts more than one `input_audio` reference clip per request. Requests carrying several clips are only routed to endpoints where this is true.
	SupportsMultipleAudioReferences *bool `json:"supports_multiple_audio_references" yaml:"supports_multiple_audio_references"`
	// Per-variant `tool_choice` support. `tool_choice` in `supported_parameters` only says the parameter is accepted; these flags say which of its values passed testing.
	// Example: {"auto": true, "function": true, "none": true, "required": true}
	SupportsToolChoice ListEndpointsResponseX200DataEndpointsItemSupportsToolChoice `json:"supports_tool_choice" yaml:"supports_tool_choice"`
	// Whether this TTS endpoint accepts inline reference audio (`input_references`) for stateless voice cloning. Requests carrying reference audio are only routed to endpoints where this is true.
	SupportsVoiceCloning *bool  `json:"supports_voice_cloning" yaml:"supports_voice_cloning"`
	Tag                  string `json:"tag" yaml:"tag"`
	// Throughput percentiles in tokens per second over the last 30 minutes. Throughput measures output token generation speed; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only visible when authenticated with an API key or cookie; returns null for unauthenticated requests.
	ThroughputLast30m *ListEndpointsResponseX200DataEndpointsItemThroughputLast30m `json:"throughput_last_30m" yaml:"throughput_last_30m"`
	// Uptime percentage over the last day: the share of minutes in which at least 80% of provider attempts succeeded, counting only minutes with 10 or more attempts. Rate-limited and caller-caused failures are excluded. Null when no minute had enough traffic.
	// format: double
	UptimeLast1d *float64 `json:"uptime_last_1d" yaml:"uptime_last_1d"`
	// format: double
	UptimeLast30m *float64 `json:"uptime_last_30m" yaml:"uptime_last_30m"`
	// Uptime percentage over the last 5 minutes: the share of minutes in which at least 80% of provider attempts succeeded, counting only minutes with 10 or more attempts. Rate-limited and caller-caused failures are excluded. Null when no minute had enough traffic.
	// format: double
	UptimeLast5m *float64 `json:"uptime_last_5m" yaml:"uptime_last_5m"`
}

// Latency percentiles in milliseconds over the last 30 minutes. Latency measures time to first token. Only visible when authenticated with an API key or cookie; returns null for unauthenticated requests.
type ListEndpointsResponseX200DataEndpointsItemLatencyLast30m struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemNativeToolsValue struct {
	// The provider tool type the request is translated to when this tool runs natively, e.g. `web_search_20260209` on Anthropic or `google_search` on Gemini.
	Type string `json:"type" yaml:"type"`
}

// Endpoint performance over the last 30 minutes, keyed by the kind of request served (e.g. `text_generation`, `image_generation`). Additive to the legacy singular latency and throughput fields; image and video generation report end-to-end latency. Only visible when authenticated with an API key or cookie.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkload struct {
	Decisions       *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadDecisions       `json:"decisions" yaml:"decisions"`
	Embeddings      *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadEmbeddings      `json:"embeddings" yaml:"embeddings"`
	ImageGeneration *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadImageGeneration `json:"image_generation" yaml:"image_generation"`
	Rerank          *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadRerank          `json:"rerank" yaml:"rerank"`
	Stt             *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadStt             `json:"stt" yaml:"stt"`
	TextGeneration  *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTextGeneration  `json:"text_generation" yaml:"text_generation"`
	Tts             *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTts             `json:"tts" yaml:"tts"`
	Unknown         *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadUnknown         `json:"unknown" yaml:"unknown"`
	VideoGeneration *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadVideoGeneration `json:"video_generation" yaml:"video_generation"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadDecisions struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadDecisionsLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadDecisionsThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadDecisionsLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadDecisionsThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadEmbeddings struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadEmbeddingsLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadEmbeddingsThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadEmbeddingsLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadEmbeddingsThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadImageGeneration struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadImageGenerationLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadImageGenerationThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadImageGenerationLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadImageGenerationThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadRerank struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadRerankLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadRerankThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadRerankLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadRerankThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadStt struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadSttLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadSttThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadSttLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadSttThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTextGeneration struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTextGenerationLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTextGenerationThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTextGenerationLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTextGenerationThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTts struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTtsLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTtsThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTtsLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadTtsThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadUnknown struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadUnknownLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadUnknownThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadUnknownLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadUnknownThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadVideoGeneration struct {
	// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
	Latency *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadVideoGenerationLatency `json:"latency" yaml:"latency"`
	// Total requests admitted for this workload in the window.
	RequestCount *int64 `json:"request_count" yaml:"request_count"`
	// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
	Throughput *ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadVideoGenerationThroughput `json:"throughput" yaml:"throughput"`
}

// Latency percentiles in milliseconds for this workload. Image and video generation report full end-to-end generation time, because their raw latency only measures request acknowledgement; every other workload reports time to first token or result.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadVideoGenerationLatency struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

// Throughput percentiles in tokens per second; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only meaningful for text generation; null for workloads without token throughput.
type ListEndpointsResponseX200DataEndpointsItemPerfLast30mByWorkloadVideoGenerationThroughput struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}

type ListEndpointsResponseX200DataEndpointsItemPricing struct {
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
	Overrides *[]ListEndpointsResponseX200DataEndpointsItemPricingOverridesItem `json:"overrides" yaml:"overrides"`
	// Price in USD per token for prompt (input) processing
	Prompt string `json:"prompt" yaml:"prompt"`
	// Price in USD per request
	Request *string `json:"request" yaml:"request"`
	// Price in USD per web search
	WebSearch *string `json:"web_search" yaml:"web_search"`
}

// A conditional override of the base pricing. An entry applies only when all of its condition fields (e.g. min_prompt_tokens, or the utc_start/utc_end time window) match the request; among applicable entries, later entries win per price key; price keys absent from an entry inherit the base price.
type ListEndpointsResponseX200DataEndpointsItemPricingOverridesItem struct {
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

// Per-variant `tool_choice` support. `tool_choice` in `supported_parameters` only says the parameter is accepted; these flags say which of its values passed testing.
type ListEndpointsResponseX200DataEndpointsItemSupportsToolChoice struct {
	// Whether the endpoint supports `"tool_choice": "auto"`.
	Auto bool `json:"auto" yaml:"auto"`
	// Whether the endpoint supports naming a function, `"tool_choice": {"type": "function", ...}`.
	Function bool `json:"function" yaml:"function"`
	// Whether the endpoint supports `"tool_choice": "none"`.
	None bool `json:"none" yaml:"none"`
	// Whether the endpoint supports `"tool_choice": "required"`.
	Required bool `json:"required" yaml:"required"`
}

// Throughput percentiles in tokens per second over the last 30 minutes. Throughput measures output token generation speed; pN is the speed that N% of requests reached or exceeded, so p90 is lower than p50. Only visible when authenticated with an API key or cookie; returns null for unauthenticated requests.
type ListEndpointsResponseX200DataEndpointsItemThroughputLast30m struct {
	// Median (50th percentile)
	// Example: 25.5
	// format: double
	P50 float64 `json:"p50" yaml:"p50"`
	// 75th percentile
	// Example: 35.2
	// format: double
	P75 float64 `json:"p75" yaml:"p75"`
	// 90th percentile
	// Example: 48.7
	// format: double
	P90 float64 `json:"p90" yaml:"p90"`
	// 99th percentile
	// Example: 85.3
	// format: double
	P99 float64 `json:"p99" yaml:"p99"`
}
