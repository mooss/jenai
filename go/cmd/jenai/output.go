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
	"gopkg.in/yaml.v3"
)

func cleanClose(f *os.File) {
	if f != nil {
		f.Close()
	}
}

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

	teeFile := openTeeFile(teePath)
	if teeFile != nil {
		defer cleanClose(teeFile)
		header, err := teeHeader(cfg, prompt, spec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't build tee header (will proceed nonetheless): %s\n", err)
			teeFile.Close()
			teeFile = nil
		} else if _, err := io.WriteString(teeFile, header); err != nil {
			fmt.Fprintf(os.Stderr, "can't tee to %s (will proceed nonetheless): %s\n", cfg.TeeFile, err)
			teeFile.Close()
			teeFile = nil
		}
	}
	thoughtlessFile := openTeeFile(thoughtlessPath)
	if thoughtlessFile != nil {
		defer cleanClose(thoughtlessFile)
	}
	var teeReasoningStarted, teeReasoningClosed bool
	teeWrite := func(f **os.File, path, s string) {
		if *f == nil {
			return
		}
		if _, err := io.WriteString(*f, s); err != nil {
			fmt.Fprintf(os.Stderr, "can't tee to %s (will proceed nonetheless): %s\n", path, err)
			(*f).Close()
			*f = nil
		}
	}

	onReasoning := func(s string) {
		if teeFile != nil {
			if !teeReasoningStarted {
				teeWrite(&teeFile, cfg.TeeFile, "<think>\n")
				teeReasoningStarted = true
			}
			teeWrite(&teeFile, cfg.TeeFile, s)
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

	onContent := func(s string) {
		if teeFile != nil {
			if teeReasoningStarted && !teeReasoningClosed {
				teeWrite(&teeFile, cfg.TeeFile, "\n</think>\n")
				teeReasoningClosed = true
			}
			teeWrite(&teeFile, cfg.TeeFile, s)
		}
		if thoughtlessFile != nil {
			teeWrite(&thoughtlessFile, cfg.ThoughtlessTeeFile, s)
		}
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

func openTeeFile(path string) *os.File {
	if path == "" {
		return nil
	}
	full := canonicalTeePath(path)
	f, err := os.OpenFile(full, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't tee to %s (will proceed nonetheless): %s\n", full, err)
		return nil
	}
	return f
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
