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

	"github.com/mooss/bagend/go/cmd"
	"github.com/mooss/bagend/go/flag"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/prompts"
)

func main() {
	if err := ensureConfig(); err != nil {
		fatal(err)
	}

	jenai := setupCommand(noerr(prompts.Embedded()))
	noerr0(jenai.Execute(os.Args[1:]))
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

func setupCommand(library prompts.Library) *cmd.Command {
	cfg := config.Jenai{}
	dir, err := config.Root()
	if err != nil {
		fatal(err)
	}
	noerr0(cfg.LoadFile(filepath.Join(dir, "config.yaml")))

	jenai := cmd.New(os.Args[0], "CLI prompt builder").Usage("COMMAND ...ARGS")
	jenai.Runner(func() error {
		fmt.Print(jenai.Help())
		return nil
	})

	promptFlags := cfg.LLMRequestFlags()
	prompt := cmd.New("prompt", "Use a prompt from the prompt library").
		Alias("p").EmplaceFlags(promptFlags).Runner(func() error {
		cfg.Positional = promptFlags.Positional
		return run(&cfg, func() (config.Prompt, error) { return cfg.BuildPrompt(library) })
	})

	oneshotFlags := cfg.LLMRequestFlags()
	oneshot := cmd.New("oneshot", "Create an ad-hoc prompt from positional arguments").
		Alias("o").EmplaceFlags(oneshotFlags).Runner(func() error {
		cfg.Positional = oneshotFlags.Positional
		return run(&cfg, cfg.BuildOneshotPrompt)
	})

	modelZoo := cmd.New("model-zoo", "List the content of the model zoo").Alias("mz").
		Runner(func() error {
			printAvailableModels()
			return nil
		})

	capaFlags := flag.NewParser().WithPositionalArguments()
	source := ""
	capaFlags.String("source", &source, "Capability source: overview, endpoints, or rawyaml").
		Default("overview").Allowed("overview", "endpoints", "rawyaml")
	capa := cmd.New("capa", "List a model's OpenRouter capabilities").
		EmplaceFlags(capaFlags).Runner(func() error {
		return runCapa(&cfg, capaFlags.Positional, source)
	})

	promptLibrary := cmd.New("prompt-library", "List the content of the prompt library").
		Alias("pl").Runner(func() error {
		printEmbeddedPromptNames(library.Prompts)
		return nil
	})

	jenai.RegisterSubs(oneshot, prompt, capa, modelZoo, promptLibrary)

	return jenai
}

func run(cfg *config.Jenai, mkprompt func() (config.Prompt, error)) error {
	if cfg.DumpConfig {
		spec := noerr(modelSpec(cfg))
		fmt.Println(pretty(map[string]any{
			"model":  spec,
			"config": cfg,
		}))
		os.Exit(0)
	}

	prompt := noerr(mkprompt())

	if cfg.DryRun {
		fmt.Println(prompt)
		os.Exit(0)
	}

	if prompt.Empty() {
		return errors.New("the prompt is empty")
	}

	spec := noerr(modelSpec(cfg))
	apiKey := noerr(cfg.APIKeyFor(spec.Provider))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return streamAndDisplay(ctx, cfg, spec, apiKey, prompt, prompt.String())
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
