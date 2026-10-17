package main

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/mooss/bagend/go/flag"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/prompts"
)

func main() {
	if err := ensureConfig(); err != nil {
		fatal(err)
	}
	cfg, parser := conf()

	if len(os.Args) == 1 {
		fmt.Print(parser.Help())
		os.Exit(0)
	}

	dumpConfig := false
	parser.Bool("dump-config", &dumpConfig, "print the config and exit")

	noerr0(cfg.ParseCLI(parser, os.Args[1:]))

	library := noerr(prompts.Embedded())

	if cfg.List {
		printEmbeddedPromptNames(library.Prompts)
		os.Exit(0)
	}

	if cfg.ListModels {
		printAvailableModels()
		os.Exit(0)
	}

	if dumpConfig {
		spec := noerr(modelSpec(cfg))
		fmt.Println(pretty(map[string]any{
			"model":  spec,
			"config": cfg,
		}))
		os.Exit(0)
	}

	run(cfg, library)
}

func printEmbeddedPromptNames(promptNames map[string]string) {
	for name := range promptNames {
		fmt.Println(name)
	}
}

func printAvailableModels() {
	specs, err := models.ModelSpecs()
	if err != nil {
		fatal(err)
	}

	longest := 0
	for _, spec := range specs {
		if len(spec.ShortName) > longest {
			longest = len(spec.ShortName)
		}
	}

	format := fmt.Sprintf("%%-%ds  (%%s/%%s)\n", longest)
	for _, short := range slices.Sorted(maps.Keys(specs)) {
		spec := specs[short]
		fmt.Printf(format, spec.ShortName, spec.Provider, spec.Author)
	}
}

//////////////////////////
// High-level utilities //

func conf() (*config.Jenai, *flag.Parser) {
	cfg := config.Jenai{}
	dir, err := config.Root()
	if err != nil {
		fatal(err)
	}
	noerr0(cfg.LoadFile(filepath.Join(dir, "config.yaml")))
	parser := cfg.RegisterCLI()
	flag.WithHelp(os.Args[0], "PROMPT ...ARGS")(parser)

	return &cfg, parser
}

func run(cfg *config.Jenai, lib prompts.Library) {
	prompt := noerr(cfg.BuildPrompt(lib))

	if cfg.DryRun {
		fmt.Println(prompt)
		os.Exit(0)
	}

	if prompt.Empty() {
		fatal(errors.New("the prompt is empty"))
	}

	spec := noerr(modelSpec(cfg))
	apiKey := noerr(cfg.APIKeyFor(spec.Provider))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := streamAndDisplay(ctx, cfg, spec, apiKey, prompt, prompt.String()); err != nil {
		fatal(err)
	}
}

func modelSpec(cfg *config.Jenai) (models.Spec, error) {
	fromZoo := func() (models.Spec, error) { return models.Get(cfg.Model) }

	// Qualified name in provider:author/model formet.
	provider, rest, found := strings.Cut(cfg.Model, ":")
	if !found {
		return fromZoo()
	}

	author, model, found := strings.Cut(rest, "/")
	if !found {
		return fromZoo()
	}

	return models.Spec{Provider: provider, Author: author, Model: model}, nil
}
