//revive:disable:line-length-limit
package openrouter

// ListEndpointsResponse is the successful response body.
type ListEndpointsResponse struct {
	Data ModelEndpointsData `json:"data" yaml:"data"`
}

type ModelEndpointsData struct {
	ID           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Created      int64             `json:"created" yaml:"created"`
	Description  string            `json:"description" yaml:"description"`
	Architecture ModelArchitecture `json:"architecture" yaml:"architecture"`
	Endpoints    []ModelEndpoint   `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
}

type ModelEndpoint struct {
	Name                            string                             `json:"name" yaml:"name"`
	ModelID                         string                             `json:"model_id" yaml:"model_id"`
	ModelName                       string                             `json:"model_name" yaml:"model_name"`
	ContextLength                   *int64                             `json:"context_length,omitempty" yaml:"context_length,omitempty"`
	Pricing                         ModelPricing                       `json:"pricing" yaml:"pricing"`
	ProviderName                    string                             `json:"provider_name" yaml:"provider_name"`
	Tag                             string                             `json:"tag" yaml:"tag"`
	Quantization                    *string                            `json:"quantization" yaml:"quantization"`
	Status                          int                                `json:"status" yaml:"status"`
	MaxCompletionTokens             *int64                             `json:"max_completion_tokens" yaml:"max_completion_tokens"`
	MaxPromptTokens                 *int64                             `json:"max_prompt_tokens" yaml:"max_prompt_tokens"`
	SupportedParameters             []string                           `json:"supported_parameters" yaml:"supported_parameters"`
	SupportsToolChoice              EndpointToolChoiceSupport          `json:"supports_tool_choice" yaml:"supports_tool_choice"`
	UptimeLast30m                   *float64                           `json:"uptime_last_30m" yaml:"uptime_last_30m"`
	UptimeLast5m                    *float64                           `json:"uptime_last_5m" yaml:"uptime_last_5m"`
	UptimeLast1d                    *float64                           `json:"uptime_last_1d" yaml:"uptime_last_1d"`
	SupportsImplicitCaching         bool                               `json:"supports_implicit_caching" yaml:"supports_implicit_caching"`
	NativeTools                     map[string]EndpointNativeTool      `json:"native_tools" yaml:"native_tools"`
	SupportsDiarization             bool                               `json:"supports_diarization" yaml:"supports_diarization"`
	LatencyLast30m                  *EndpointPercentiles               `json:"latency_last_30m" yaml:"latency_last_30m"`
	ThroughputLast30m               *EndpointPercentiles               `json:"throughput_last_30m" yaml:"throughput_last_30m"`
	PerfLast30mByWorkload           map[string]EndpointWorkloadMetrics `json:"perf_last_30m_by_workload,omitempty" yaml:"perf_last_30m_by_workload,omitempty"`
	SupportsImageReference          *bool                              `json:"supports_image_reference,omitempty" yaml:"supports_image_reference,omitempty"`
	SupportsMultipleAudioReferences *bool                              `json:"supports_multiple_audio_references,omitempty" yaml:"supports_multiple_audio_references,omitempty"`
	SupportsVoiceCloning            *bool                              `json:"supports_voice_cloning,omitempty" yaml:"supports_voice_cloning,omitempty"`
}

type EndpointPercentiles struct {
	P50 float64 `json:"p50" yaml:"p50"`
	P75 float64 `json:"p75" yaml:"p75"`
	P90 float64 `json:"p90" yaml:"p90"`
	P99 float64 `json:"p99" yaml:"p99"`
}

type EndpointWorkloadMetrics struct {
	Latency      *EndpointPercentiles `json:"latency" yaml:"latency"`
	Throughput   *EndpointPercentiles `json:"throughput" yaml:"throughput"`
	RequestCount *int64               `json:"request_count" yaml:"request_count"`
}

type EndpointNativeTool struct {
	Type string `json:"type" yaml:"type"`
}

type EndpointToolChoiceSupport struct {
	Auto     bool `json:"auto" yaml:"auto"`
	Function bool `json:"function" yaml:"function"`
	None     bool `json:"none" yaml:"none"`
	Required bool `json:"required" yaml:"required"`
}
