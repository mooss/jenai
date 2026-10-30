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

	"github.com/mooss/bagend/go/sysu"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/highlight"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/utils"
	"gopkg.in/yaml.v3"
)

func streamAndDisplay(
	ctx context.Context, qry *config.Query, spec models.Spec, apiKey string, prompt config.Prompt,
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

	out := NewTerminalRenderer(qry, teeFile, thoughtlessFile)
	streamer, err := openrouter.NewStreamer(ctx, spec, apiKey, prompt.String(), out)
	if err != nil {
		return "", err
	}

	defer func() {
		if usage := streamer.FormatUsage(); usage != "" {
			fmt.Fprint(os.Stderr, "\n", usage)
		}
	}()

	if err := streamer.ProcessChunks(); err != nil {
		return "", err
	}

	if out.UnclosedReasoning() {
		return "", errors.New("reasoning never closed")
	}
	if err := out.Finish(); err != nil {
		return "", err
	}

	return out.Response(), nil
}

// renderAndDisplayMarkdown renders Markdown content with syntax highlighting and displays it.
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

//////////////////////
// CompletionReader //

type TerminalRenderer struct {
	qry             *config.Query
	output          *highlight.MarkdownStreamer
	reasoningOutput *highlight.MarkdownStreamer
	teeFile         *utils.TeeFile
	thoughtlessFile *utils.TeeFile
	stdoutThink     thinkGuard
	teeThink        thinkGuard
	endsWithNewline bool
	response        strings.Builder
}

func NewTerminalRenderer(
	qry *config.Query, teeFile, thoughtlessFile *utils.TeeFile,
) *TerminalRenderer {
	color := sysu.TerminalColorEnabled()
	return &TerminalRenderer{
		qry:             qry,
		output:          highlight.NewMarkdownStreamer(os.Stdout, color),
		reasoningOutput: highlight.NewMarkdownStreamer(os.Stdout, color),
		teeFile:         teeFile,
		thoughtlessFile: thoughtlessFile,
		endsWithNewline: true, // True by default to avoid always calling onContent.
	}
}

func (r *TerminalRenderer) ReadReasoning(s string) error {
	if r.teeFile.IsOpen() {
		r.teeThink.open(func() { r.teeFile.Write("<think>\n") })
		r.teeFile.Write(s)
	}
	if r.qry.Thoughtless {
		return nil
	}
	var openErr error
	r.stdoutThink.open(func() {
		tag := r.reasoningOutput.ColorizeTag("<think>")
		if _, err := fmt.Fprint(os.Stdout, tag, "\n"); err != nil {
			openErr = err
		}
	})
	if openErr != nil {
		return openErr
	}
	return r.reasoningOutput.Write(s)
}

func (r *TerminalRenderer) ReadContent(s string) error {
	r.response.WriteString(s)
	return r.writeContent(s)
}

func (r *TerminalRenderer) writeContent(s string) error {
	if r.teeFile.IsOpen() {
		r.teeThink.close(func() { r.teeFile.Write("\n</think>\n") })
		r.teeFile.Write(s)
	}
	if r.thoughtlessFile.IsOpen() {
		r.thoughtlessFile.Write(s)
	}
	if !r.qry.Thoughtless {
		var closeErr error
		r.stdoutThink.close(func() {
			if err := r.reasoningOutput.Flush(); err != nil {
				closeErr = err
				return
			}
			tag := r.reasoningOutput.ColorizeTag("</think>")
			if _, err := fmt.Fprint(os.Stdout, "\n", tag, "\n\n"); err != nil {
				closeErr = err
			}
		})
		if closeErr != nil {
			return closeErr
		}
	}
	r.endsWithNewline = strings.HasSuffix(s, "\n")
	return r.output.Write(s)
}

func (r *TerminalRenderer) Response() string {
	return r.response.String()
}

func (r *TerminalRenderer) UnclosedReasoning() bool {
	return r.stdoutThink.unclosed()
}

func (r *TerminalRenderer) Finish() error {
	if !r.endsWithNewline {
		if err := r.writeContent("\n"); err != nil {
			return err
		}
	}
	return r.output.Flush()
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
