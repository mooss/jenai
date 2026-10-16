package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/highlight"
	"github.com/mooss/jenai/go/llm"
	"github.com/mooss/jenai/go/models"
	"gopkg.in/yaml.v3"
)

func streamAndDisplay(ctx context.Context, cfg *config.Jenai, spec models.Spec, apiKey, promptStr string) (content, full string) {
	var contentBuilder, reasoningBuilder strings.Builder
	output := highlight.NewMarkdownStreamer(os.Stdout, "pygments", terminalColorEnabled(os.Stdout))
	reasoningOutput := highlight.NewMarkdownStreamer(os.Stdout, "pygments", terminalColorEnabled(os.Stdout))
	var reasoningStarted, reasoningClosed bool

	onReasoning := func(s string) {
		reasoningBuilder.WriteString(s)
		if cfg.Thoughtless {
			return
		}
		if !reasoningStarted {
			fmt.Fprint(os.Stdout, reasoningOutput.ColorizeTag("<think>"), "\n")
			reasoningStarted = true
		}
		noerr0(reasoningOutput.Write(s))
	}

	onContent := func(s string) {
		contentBuilder.WriteString(s)
		if reasoningStarted && !reasoningClosed && !cfg.Thoughtless {
			noerr0(reasoningOutput.Flush())
			fmt.Fprint(os.Stdout, "\n", reasoningOutput.ColorizeTag("</think>"), "\n\n")
			reasoningClosed = true
		}
		noerr0(output.Write(s))
	}

	chunks, errs, err := llm.Stream(ctx, spec, apiKey, promptStr)
	noerr0(err)
	for chunk := range chunks {
		if len(chunk.Choices) == 0 {
			if chunk.Usage != nil {
				continue // Expected final usage-only chunk from `IncludeUsage: true`.
			}
			fmt.Fprintln(os.Stderr, "warning: chunk with empty choices and no usage")
			continue
		}

		c := chunk.Choices[0]
		if c.Delta.Reasoning != "" {
			onReasoning(c.Delta.Reasoning)
		}
		if c.Delta.Content != "" {
			onContent(c.Delta.Content)
		}

	}
	noerr0(<-errs)

	if reasoningStarted && !reasoningClosed {
		noerr0(errors.New("reasoning never closed"))
	}

	content = contentBuilder.String()
	reasoning := reasoningBuilder.String()

	noerr0(output.Flush())

	full = content
	if reasoning != "" {
		full = "<think>\n" + reasoning + "\n</think>\n" + content
	}

	if content != "" && !strings.HasSuffix(content, "\n") {
		fmt.Fprintln(os.Stdout)
	}

	content = withTrailingNewline(content)
	full = withTrailingNewline(full)

	return content, full
}

func terminalColorEnabled(writer io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func withTrailingNewline(s string) string {
	if s != "" && !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}

func writeTees(cfg *config.Jenai, prompt config.Prompt, full, content string, spec models.Spec) {
	if cfg.TeeFile != "" {
		if err := tee(cfg, prompt, full, spec); err != nil {
			fmt.Fprintf(os.Stderr,
				"can't tee to %s (will proceed nonetheless): %s", cfg.TeeFile, err)
		}
	}

	if cfg.ThoughtlessTeeFile != "" {
		if err := thoughtlessTee(cfg, content); err != nil {
			fmt.Fprintf(os.Stderr,
				"can't tee to %s (will proceed nonetheless): %s", cfg.ThoughtlessTeeFile, err)
		}
	}
}

func tee(cfg *config.Jenai, prompt config.Prompt, answer string, spec models.Spec) error {
	metadata, err := structToMap(cfg)
	if err != nil {
		return err
	}

	metadata["date"] = time.Now().Format("2006-01-02")
	metadata["model"] = spec.QualifiedName()
	metadata["prompt"] = prompt.Static()
	delete(metadata, "positional") // Noise that should be present in `prompt` anyway.

	yamlmeta, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	content := strings.Join([]string{"---", string(yamlmeta) + "---\n", answer}, "\n")
	return os.WriteFile(cfg.TeeFile, []byte(content), 0644)
}

func thoughtlessTee(cfg *config.Jenai, answer string) error {
	return os.WriteFile(cfg.ThoughtlessTeeFile, []byte(answer), 0644)
}
