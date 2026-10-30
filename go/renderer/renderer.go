// Package renderer handles rendering of LLM completion.
package renderer

type CompletionReader interface {
	ReadReasoning(string) error
	ReadContent(string) error
}
