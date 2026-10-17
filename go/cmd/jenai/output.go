package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/highlight"
	"github.com/mooss/jenai/go/llm"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/utils"
	"gopkg.in/yaml.v3"
)

func streamAndDisplay(
	ctx context.Context, cfg *config.Jenai, spec models.Spec, apiKey string,
	prompt config.Prompt, promptStr string,
) (err error) {
	teePath := canonicalTeePath(cfg.TeeFile)
	thoughtlessPath := canonicalTeePath(cfg.ThoughtlessTeeFile)
	if teePath != "" && teePath == thoughtlessPath {
		return fmt.Errorf("tee file and thoughtless tee file are the same file: %s", teePath)
	}

	output := highlight.NewMarkdownStreamer(os.Stdout, "pygments", terminalColorEnabled(os.Stdout))
	reasoningOutput := highlight.NewMarkdownStreamer(os.Stdout, "pygments", terminalColorEnabled(os.Stdout))
	var reasoningStarted, reasoningClosed bool

	teeFile := utils.OpenTeeFile(teePath)
	if teeFile.IsOpen() {
		defer teeFile.Finish()
		header, err := teeHeader(cfg, prompt, spec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't build tee header (will proceed nonetheless): %s\n", err)
			teeFile.Finish()
		} else {
			teeFile.Write(header)
		}
	}
	thoughtlessFile := utils.OpenTeeFile(thoughtlessPath)
	if thoughtlessFile.IsOpen() {
		defer thoughtlessFile.Finish()
	}
	var teeReasoningStarted, teeReasoningClosed bool

	onReasoning := func(s string) {
		if teeFile.IsOpen() {
			if !teeReasoningStarted {
				teeFile.Write("<think>\n")
				teeReasoningStarted = true
			}
			teeFile.Write(s)
		}
		if cfg.Thoughtless {
			return
		}
		if !reasoningStarted {
			fmt.Fprint(os.Stdout, reasoningOutput.ColorizeTag("<think>"), "\n")
			reasoningStarted = true
		}
		noerr0(reasoningOutput.Write(s))
	}

	endsWithNewline := true // True by default to avoid always calling onContent.
	onContent := func(s string) {
		if teeFile.IsOpen() {
			if teeReasoningStarted && !teeReasoningClosed {
				teeFile.Write("\n</think>\n")
				teeReasoningClosed = true
			}
			teeFile.Write(s)
		}
		if thoughtlessFile.IsOpen() {
			thoughtlessFile.Write(s)
		}
		if reasoningStarted && !reasoningClosed && !cfg.Thoughtless {
			noerr0(reasoningOutput.Flush())
			fmt.Fprint(os.Stdout, "\n", reasoningOutput.ColorizeTag("</think>"), "\n\n")
			reasoningClosed = true
		}

		endsWithNewline = strings.HasSuffix(s, "\n")
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

	if !endsWithNewline {
		onContent("\n")
	}
	noerr0(output.Flush())

	return nil
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

func canonicalTeePath(path string) string {
	if path == "" {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func teeHeader(cfg *config.Jenai, prompt config.Prompt, spec models.Spec) (string, error) {
	metadata, err := structToMap(cfg)
	if err != nil {
		return "", err
	}

	metadata["date"] = time.Now().Format("2006-01-02")
	metadata["model"] = spec.QualifiedName()
	metadata["prompt"] = prompt.Static()
	delete(metadata, "positional") // Noise that should be present in `prompt` anyway.

	yamlmeta, err := yaml.Marshal(metadata)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{"---", string(yamlmeta) + "---\n"}, "\n") + "\n", nil
}
