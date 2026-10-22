// Package config defines jenai's configuration.
package config

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

	"github.com/mooss/bagend/go/flag"
	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/jenai/go/prompts"
	"github.com/mooss/jenai/go/utils"
)

func Root() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "jenai"), nil
}

//revive:disable:line-length-limit
type Query struct {
	APIKeys            map[string]string `yaml:"-" json:"-"`
	Context            Context           `yaml:"context" json:"context"`
	DryRun             bool              `yaml:"dryrun,omitempty" json:"dryrun,omitempty"`
	Model              string            `yaml:"model,omitempty" json:"model,omitempty"`
	Paste              bool              `yaml:"paste,omitempty" json:"paste,omitempty"`
	Positional         []string          `yaml:"positional,omitempty" json:"positional,omitempty"`
	Reasoning          string            `yaml:"reasoning,omitempty" json:"reasoning,omitempty"`
	Thoughtless        bool              `yaml:"thoughtless,omitempty" json:"thoughtless,omitempty"`
	ThoughtlessTeeFile string            `yaml:"thoughtlessteefile,omitempty" json:"thoughtlessteefile,omitempty"`
	TeeFile            string            `yaml:"teefile,omitempty" json:"teefile,omitempty"`

	DumpConfig bool `yaml:"-" json:"-"`
}

/////////////////////////////////
// Construction and validation //

func (qry *Query) BindFlags() *flag.Parser {
	flags := flag.NewParser()

	flags.Positional(&qry.Positional)
	flags.Bool("context-above", &qry.Context.Above,
		"Put context files and dir above instructions")
	flags.StringSlice("dir", &qry.Context.Dirs, "Include all files in directory as context")
	flags.Bool("dry-run", &qry.DryRun, "Print interpolated prompt without sending to LLM").
		Alias("n")
	flags.StringSlice("file", &qry.Context.Files, "Include specific file(s) as context")
	flags.Bool("linum", &qry.Context.LineNumbers, "Print files with line numbers")
	flags.String(
		"model", &qry.Model, "Model name (short name from `model-zoo` or provider:author/model)",
	).Alias("m").Default("spark")
	flags.Bool("paste", &qry.Paste, "Use clipboard content as prompt")
	flags.String("reasoning", &qry.Reasoning, "Reasoning effort (must be supported by the model)")
	flags.String("tee", &qry.TeeFile, "Output first answer to both stdout and FILE (overwritten)")
	flags.Bool("thoughtless", &qry.Thoughtless, "Hide model reasoning").Alias("t")
	flags.String("thoughtless-tee", &qry.ThoughtlessTeeFile,
		"Output first answer to FILE without the reasoning block (overwritten)").
		Alias("tt")

	flags.Bool("dump-config", &qry.DumpConfig, "dump the config to stdout and exit")

	return flags
}

func (qry *Query) LoadFile(path string) error {
	loaded, err := utils.FromYAMLFile[struct {
		APIKeys map[string]string `yaml:"api-keys" json:"api-keys"`
	}](path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to load config file %q: %w", path, err)
	}

	if loaded.APIKeys != nil {
		qry.APIKeys = loaded.APIKeys
	}
	return nil
}

func (qry *Query) APIKeyFor(provider string) (string, error) {
	if key := qry.APIKeys[provider]; key != "" {
		return key, nil
	}
	dir, err := Root()
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf(
		"missing api key for %q: put it in %s as:\napi-keys:\n  %s: API_KEY_HERE",
		provider, filepath.Join(dir, "config.yaml"), provider,
	)
}

// BuildPrompt returns the complete prompt, taking all sources into account (prompt, clipboard,
// positional argument, and stdin).
// The prompt is evaluated.
func (qry *Query) BuildPrompt(lib prompts.Library) (Prompt, error) {
	primary, positional, err := qry.buildPositional(lib)
	if err != nil {
		return Prompt{}, err
	}

	return qry.buildPrompt(primary, positional)
}

// BuildOneshotPrompt returns the complete prompt for an ad-hoc positional prompt.
func (qry *Query) BuildOneshotPrompt() (Prompt, error) {
	return qry.buildPrompt("", qry.Positional)
}

func (qry *Query) buildPrompt(primary string, positional []string) (Prompt, error) {
	clipboard, err := qry.buildClipboard()
	if err != nil {
		return Prompt{}, err
	}

	contextText, paths, err := qry.Context.Build()
	if err != nil {
		return Prompt{}, err
	}

	stdin, err := readStdin()
	if err != nil {
		return Prompt{}, err
	}

	return Prompt{
		Clipboard:    clipboard,
		Context:      contextText,
		ContextAbove: qry.Context.Above,
		Paths:        paths,
		Positional:   strings.Join(positional, " "),
		Primary:      primary,
		Stdin:        stdin,
	}, nil
}

func (qry *Query) buildPositional(lib prompts.Library) (string, []string, error) {
	positional := qry.Positional
	if len(positional) == 0 {
		return "", positional, nil
	}

	positional = positional[1:]
	primary, err := prompts.NewEvalContext(lib, &positional).Evaluate(qry.Positional[0])
	if err != nil {
		return "", nil, err
	}

	return primary, positional, nil
}

func (qry *Query) buildClipboard() (string, error) {
	if !qry.Paste {
		return "", nil
	}

	return readClipboard()
}

///////////////////////
// Utility functions //

// isWayland reports whether the current session looks like Wayland.
func isWayland() bool {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	return strings.EqualFold(os.Getenv("XDG_SESSION_TYPE"), "wayland")
}

// clipboardCommands returns clipboard reader commands in preferred order.
// Wayland sessions prefer wl-paste, X11 sessions prefer xclip, with fallback to the other in case
// of misconfigured environment or missing binaries.
func clipboardCommands() [][]string {
	xclip := []string{"xclip", "-o", "-selection", "clipboard"}
	xsel := []string{"xsel", "--clipboard", "--output"}
	wlPaste := []string{"wl-paste"}

	if isWayland() {
		return [][]string{wlPaste, xclip, xsel}
	}
	return [][]string{xclip, xsel, wlPaste}
}

// readClipboard returns the content of the clipboard.
func readClipboard() (string, error) {
	var tried []string
	var errs []error

	clipboardCmds := clipboardCommands()
	for _, args := range clipboardCmds {
		if _, err := exec.LookPath(args[0]); err != nil {
			continue
		}
		tried = append(tried, args[0])

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		output, err := exec.CommandContext(ctx, args[0], args[1:]...).Output() //nolint:gosec
		cancel()
		if err == nil {
			return string(output), nil
		} // Tool exists but failed (e.g. empty clipboard), simply try the next one.
		errs = append(errs, fmt.Errorf("%s: %w", args[0], err))
	}

	if len(tried) == 0 {
		cmds := lie.Map(func(el []string) string { return el[0] }, clipboardCmds)
		return "", fmt.Errorf("failed to get clipboard content: no clipboard tool found (tried %s)",
			strings.Join(cmds, ","))
	}
	return "", fmt.Errorf("failed to get clipboard content (tried %s): %w",
		strings.Join(tried, ", "), errors.Join(errs...))
}

// readStdin returns the content of stdin, if something was piped in.
func readStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil // No data being piped.
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}
	return string(data), nil
}
