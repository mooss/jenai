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
	"github.com/mooss/bagend/go/fun"
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

func printAvailableModels() error {
	specs, err := models.ModelSpecs()
	if err != nil {
		return err
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

	return nil
}

//////////////////////////
// High-level utilities //

//revive:disable:function-length
func setupCommand(library prompts.Library) *cmd.Command {
	qry := config.Query{}
	dir, err := config.Root()
	if err != nil {
		fatal(err)
	}
	noerr0(qry.LoadFile(filepath.Join(dir, "config.yaml")))

	jenai := &cmd.Command{
		Name:    os.Args[0],
		Docline: "CLI prompt builder",
		Usage:   "COMMAND ARGS...",
	}
	jenai.Runner = func() error {
		fmt.Print(jenai.Help())
		return nil
	}

	oneshot := &cmd.Command{
		Name:    "oneshot",
		Docline: "Create an ad-hoc prompt from arguments, clipboard or stdin",
		Usage:   "[INSTRUCTIONS...]",
		Aliases: []string{"o"},
		Flags:   qry.BindFlags(),
	}
	oneshot.Runner = func() error {
		qry.Positional = oneshot.Flags.Positional
		return run(&qry, qry.BuildOneshotPrompt)
	}

	prompt := &cmd.Command{
		Name:    "prompt",
		Docline: "Use a prompt from the prompt library",
		Usage:   "NAME [ARGS...]",
		Aliases: []string{"p"},
		Flags:   qry.BindFlags(),
	}
	prompt.Runner = func() error {
		qry.Positional = prompt.Flags.Positional
		return run(&qry, func() (config.Prompt, error) { return qry.BuildPrompt(library) })
	}

	writeArg := []string{}
	force := false
	edit := &cmd.Command{
		Name:    "edit",
		Docline: "Edit files according to instructions",
		Usage:   "[INSTRUCTIONS...] --write FILE...",
		Aliases: []string{"ed"},
		Flags:   qry.BindFlags(),
	}
	edit.Flags.StringSlice("write", &writeArg, "Files to edit")
	edit.Flags.Bool("force", &force, "Allow editing files with unstaged changes")
	edit.Runner = func() error {
		qry.Positional = edit.Flags.Positional
		return runEdit(&qry, writeArg, force)
	}

	modelZoo := &cmd.Command{
		Name:    "model-zoo",
		Docline: "Show configured model shortcuts",
		Aliases: []string{"mz"},
		Runner:  printAvailableModels,
	}

	details := &cmd.Command{
		Name:    "details",
		Docline: "Print details for a model, or the whole model zoo",
		Usage:   "[MODEL]",
		Flags:   flag.NewParser().WithPositionalArguments(),
	}
	details.Runner = func() error { return printModelDetails(details.Flags.Positional) }

	regen := &cmd.Command{
		Name:        "regen",
		Docline:     "Overwrites the given model zoo YAML with current OpenRouter info",
		Usage:       "OUTPUT-FILE",
		Description: `OUTPUT-FILE will double as the model zoo source if valid.`,
		Flags:       flag.NewParser().WithPositionalArguments(),
	}
	regen.Runner = func() error {
		return regenModelZoo(&qry, regen.Flags.Positional)
	}

	modelZoo.RegisterSubs(details, regen)

	source := ""
	capa := &cmd.Command{
		Name:    "capa",
		Docline: "List a model's OpenRouter capabilities",
		Usage:   "MODEL",
		Flags:   flag.NewParser().WithPositionalArguments(),
	}
	capa.Flags.String(
		"source", &source, "Capability source: overview, endpointless, endpoints, or rawyaml",
	).Default("overview").Allowed("overview", "endpointless", "endpoints", "rawyaml")
	capa.Runner = func() error { return runCapa(&qry, capa.Flags.Positional, source) }

	promptLibrary := &cmd.Command{
		Name:    "prompt-library",
		Docline: "List the content of the prompt library",
		Aliases: []string{"pl"},
		Runner: func() error {
			printEmbeddedPromptNames(library.Prompts)
			return nil
		},
	}

	jenai.RegisterSubs(oneshot, prompt, edit, modelZoo, capa, promptLibrary)

	return jenai
}

func run(qry *config.Query, mkprompt func() (config.Prompt, error)) error {
	return runWithHooks(qry, mkprompt, nil, nil)
}

func runWithHooks(
	qry *config.Query,
	mkprompt func() (config.Prompt, error),
	prepare func() error,
	onResponse func(string) error,
) error {
	if qry.DumpConfig {
		spec := noerr(modelSpec(qry))
		prettyPrint(map[string]any{
			"model":  spec,
			"config": qry,
		})
		os.Exit(0)
	}

	if prepare != nil {
		if err := prepare(); err != nil {
			return err
		}
	}

	prompt, err := mkprompt()
	if err != nil {
		return err
	}

	if qry.DryRun {
		fmt.Println(prompt)
		os.Exit(0)
	}

	if prompt.Empty() {
		return errors.New("the prompt is empty")
	}

	spec, apiKey, err := modelAndAPIKey(qry)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	response, err := streamAndDisplay(ctx, qry, spec, apiKey, prompt, prompt.String())
	if err != nil {
		return err
	}
	if onResponse != nil {
		return onResponse(response)
	}
	return nil
}

// modelAndAPIKey resolves the model spec and API key for the query's model, applying the reasoning
// override when requested.
func modelAndAPIKey(qry *config.Query) (models.Spec, string, error) {
	spec, err := modelSpec(qry)
	if err != nil {
		return models.Spec{}, "", err
	}
	if qry.Reasoning != "" {
		if !slices.Contains(fun.ZeroDeref(spec.Reasoning).SupportedEfforts, qry.Reasoning) {
			msg := "reasoning effort %q is not supported by model %q"
			return models.Spec{}, "", fmt.Errorf(msg, qry.Reasoning, spec.QualifiedName())
		}

		supportedEfforts := append([]string(nil), spec.Reasoning.SupportedEfforts...)
		spec.Reasoning = &models.ReasoningSpec{
			SupportedEfforts: supportedEfforts,
			DefaultEffort:    qry.Reasoning,
		}
	}
	apiKey := noerr(qry.APIKeyFor(spec.Provider))
	return spec, apiKey, nil
}

func modelSpec(qry *config.Query) (models.Spec, error) {
	fromZoo := func() (models.Spec, error) { return models.Get(qry.Model) }

	// Qualified name in provider:author/model formet.
	provider, rest, found := strings.Cut(qry.Model, ":")
	if !found {
		return fromZoo()
	}

	author, model, found := strings.Cut(rest, "/")
	if !found {
		return fromZoo()
	}

	return models.Spec{Provider: provider, Author: author, Model: model}, nil
}
