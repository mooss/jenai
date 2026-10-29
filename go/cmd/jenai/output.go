package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mooss/bagend/go/obj/osl"
	"github.com/mooss/bagend/go/sysu"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/highlight"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/utils"
	"gopkg.in/yaml.v3"
)

func streamAndDisplay(
	ctx context.Context, qry *config.Query, spec models.Spec, apiKey string,
	prompt config.Prompt, promptStr string,
) (string, error) {
	teePath := canonicalTeePath(qry.TeeFile)
	thoughtlessPath := canonicalTeePath(qry.ThoughtlessTeeFile)
	if err := checkTeePaths(teePath, thoughtlessPath); err != nil {
		return "", err
	}

	teeFile := utils.OpenTeeFile(teePath)
	defer teeFile.Finish()
	writeTeeHeader(teeFile, qry, prompt, spec)
	thoughtlessFile := utils.OpenTeeFile(thoughtlessPath)
	defer thoughtlessFile.Finish()

	d := newDisplay(qry, teeFile, thoughtlessFile)
	defer d.printCost()

	chunks, errs, err := openrouter.StreamCompletion(ctx, spec, apiKey, promptStr)
	noerr0(err)

	for chunk := range chunks {
		d.handleChunk(chunk)
	}
	streamErr := <-errs

	noerr0(streamErr)

	if d.stdoutThink.unclosed() {
		noerr0(errors.New("reasoning never closed"))
	}
	d.finish()

	return d.response.String(), nil
}

/////////////
// display //

type display struct {
	qry             *config.Query
	output          *highlight.MarkdownStreamer
	reasoningOutput *highlight.MarkdownStreamer
	teeFile         *utils.TeeFile
	thoughtlessFile *utils.TeeFile
	stdoutThink     thinkGuard
	teeThink        thinkGuard
	endsWithNewline bool
	response        strings.Builder
	usage           *openrouter.Usage
}

func newDisplay(qry *config.Query, teeFile, thoughtlessFile *utils.TeeFile) *display {
	color := sysu.TerminalColorEnabled()
	return &display{
		qry:             qry,
		output:          highlight.NewMarkdownStreamer(os.Stdout, color),
		reasoningOutput: highlight.NewMarkdownStreamer(os.Stdout, color),
		teeFile:         teeFile,
		thoughtlessFile: thoughtlessFile,
		endsWithNewline: true, // True by default to avoid always calling onContent.
	}
}

// renderAndDisplayMarkdown renders Markdown content with syntax highlighting and displays it.
//
// If both stdout are connected to a terminal, the rendered content is streamed to less.
func renderAndDisplayMarkdown(content string) error {
	render := func(writer io.Writer) error {
		output := highlight.NewMarkdownStreamer(writer, sysu.TerminalColorEnabled())
		if err := output.Write(content); err != nil {
			return err
		}
		return output.Flush()
	}

	if !sysu.IsTerminal(os.Stdout) {
		return render(os.Stdout)
	}

	pager := exec.CommandContext(context.Background(), "less", "-R")
	pager.Stderr = os.Stderr

	pagerStdin, err := pager.StdinPipe()
	if err != nil {
		return render(os.Stdout)
	}

	pager.Stdout = os.Stdout
	if err := pager.Start(); err != nil {
		return render(os.Stdout)
	}

	return errors.Join(render(pagerStdin), pagerStdin.Close(), pager.Wait())
}

func (d *display) onReasoning(s string) {
	if d.teeFile.IsOpen() {
		d.teeThink.open(func() { d.teeFile.Write("<think>\n") })
		d.teeFile.Write(s)
	}
	if d.qry.Thoughtless {
		return
	}
	d.stdoutThink.open(func() {
		fmt.Fprint(os.Stdout, d.reasoningOutput.ColorizeTag("<think>"), "\n")
	})
	noerr0(d.reasoningOutput.Write(s))
}

func (d *display) onContent(s string) {
	if d.teeFile.IsOpen() {
		d.teeThink.close(func() { d.teeFile.Write("\n</think>\n") })
		d.teeFile.Write(s)
	}
	if d.thoughtlessFile.IsOpen() {
		d.thoughtlessFile.Write(s)
	}
	if !d.qry.Thoughtless {
		d.stdoutThink.close(func() {
			noerr0(d.reasoningOutput.Flush())
			fmt.Fprint(os.Stdout, "\n", d.reasoningOutput.ColorizeTag("</think>"), "\n\n")
		})
	}
	d.endsWithNewline = strings.HasSuffix(s, "\n")
	noerr0(d.output.Write(s))
}

func (d *display) handleChunk(chunk openrouter.Chunk) {
	if chunk.Usage != nil {
		d.usage = chunk.Usage
	}
	if len(chunk.Choices) == 0 {
		if chunk.Usage != nil {
			return // Expected final usage-only chunk from `IncludeUsage: true`.
		}
		fmt.Fprintln(os.Stderr, "warning: chunk with empty choices and no usage")
		return
	}
	c := chunk.Choices[0]
	if c.Delta.Reasoning != "" {
		d.onReasoning(c.Delta.Reasoning)
	}
	if c.Delta.Content != "" {
		d.response.WriteString(c.Delta.Content)
		d.onContent(c.Delta.Content)
	}
}

func (d *display) finish() {
	if !d.endsWithNewline {
		d.onContent("\n")
	}
	noerr0(d.output.Flush())
}

func (d *display) printCost() {
	if d.usage == nil {
		return
	}
	u := d.usage
	costStr := "¢? (unavailable)"
	if u.Cost != nil {
		costStr = fmt.Sprintf("¢%.4f", 100*(*u.Cost))
	}
	parts := osl.New(
		costStr,
		fmt.Sprintf(":prompt %d", u.PromptTokens),
		fmt.Sprintf(":completion %d", u.CompletionTokens),
		fmt.Sprintf(":total %d", u.TotalTokens),
	)
	if u.PromptTokensDetails != nil && u.PromptTokensDetails.CachedTokens > 0 {
		parts.Push(fmt.Sprintf(":cached %d", u.PromptTokensDetails.CachedTokens))
	}
	if u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0 {
		parts.Push(fmt.Sprintf(":reasoning %d", u.CompletionTokensDetails.ReasoningTokens))
	}

	fmt.Fprintf(os.Stderr, "\nCost: %s\n", strings.Join(parts, " "))
}

////////////////
// thinkGuard //

type thinkGuard struct {
	started, closed bool
}

func (g *thinkGuard) open(openFn func()) {
	if !g.started {
		openFn()
		g.started = true
	}
}

func (g *thinkGuard) close(closeFn func()) {
	if g.started && !g.closed {
		closeFn()
		g.closed = true
	}
}

func (g *thinkGuard) unclosed() bool {
	return g.started && !g.closed
}

///////////////////////
// Utility functions //

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

func checkTeePaths(teePath, thoughtlessPath string) error {
	if teePath != "" && teePath == thoughtlessPath {
		return fmt.Errorf("tee file and thoughtless tee file are the same file: %s", teePath)
	}
	return nil
}

func writeTeeHeader(
	teeFile *utils.TeeFile, qry *config.Query, prompt config.Prompt, spec models.Spec,
) {
	if !teeFile.IsOpen() {
		return
	}

	header, err := teeHeader(qry, prompt, spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't build tee header (will proceed nonetheless): %s\n", err)
		teeFile.Finish()
	} else {
		teeFile.Write(header)
	}
}

func teeHeader(qry *config.Query, prompt config.Prompt, spec models.Spec) (string, error) {
	metadata, err := structToMap(qry)
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
