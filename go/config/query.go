package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mooss/bagend/go/flag"
	"github.com/mooss/bagend/go/sysu"
	"github.com/mooss/jenai/go/prompts"
	"github.com/mooss/jenai/go/utils"
)

//revive:disable:line-length-limit
type Query struct {
	APIKeys            map[string]string `yaml:"-" json:"-"`
	Context            Context           `yaml:"context" json:"context"`
	Copy               bool              `yaml:"copy,omitempty" json:"copy,omitempty"`
	DryRun             bool              `yaml:"dryrun,omitempty" json:"dryrun,omitempty"`
	Model              string            `yaml:"model,omitempty" json:"model,omitempty"`
	Paste              bool              `yaml:"paste,omitempty" json:"paste,omitempty"`
	Positional         []string          `yaml:"positional,omitempty" json:"positional,omitempty"`
	Effort             string            `yaml:"effort,omitempty" json:"effort,omitempty"`
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
	flags.StringSlice("read", &qry.Context.Readable.Files, "Files to read")
	flags.StringSlice("read-dir", &qry.Context.Readable.Dirs, "Dirs to read (recursively)")
	flags.Bool("dry-run", &qry.DryRun, "Print interpolated prompt without sending to LLM").
		Alias("n")
	flags.Bool("linum", &qry.Context.LineNumbers, "Print files with line numbers")
	flags.String(
		"model", &qry.Model, "Model name (short name from `model-zoo` or provider:author/model)",
	).Alias("m").Default("spark")
	flags.Bool("paste", &qry.Paste, "Use clipboard content as prompt")
	flags.Bool("copy", &qry.Copy, "Copy response to clipboard (copies prompt with --dry-run)")
	flags.String("effort", &qry.Effort, "Reasoning effort (must be supported by the model)")
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

/////////////////////
// Utility methods //

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

func (qry *Query) MaybeCopy(content string) {
	if qry.Copy {
		if err := sysu.WriteClipboard(content); err != nil {
			fmt.Fprintf(os.Stderr, "clipboard writing attempt returned an error: %s\n", err)
		}
	}
}

/////////////////////
// Prompt building //

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

	stdin, err := sysu.ReadStdin()
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

	return sysu.ReadClipboard()
}
