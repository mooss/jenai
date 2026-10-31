// Package renderer handles rendering of LLM completion.
package renderer

type CompletionReader interface {
	// ReadReasoning consumes a chunk of the model's reasoning.
	ReadReasoning(string) error
	// ReadContent consumes a chunk of the completion content.
	ReadContent(string) error
}

type thinkState int

const (
	pending thinkState = iota
	thinking
	over
)
