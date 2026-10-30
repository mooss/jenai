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
	"syscall"

	"github.com/mooss/bagend/go/cmd"
	"github.com/mooss/bagend/go/flag"
	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/prompts"
)

func main() {
	if err := config.EnsureProperSetup(); err != nil {
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
		Flags:   flag.NewParser(),
	}
	jenai.Runner = func() error {
		fmt.Print(jenai.Help())
		return nil
	}
	jenai.Flags.String(
		"fmt", &config.Global.OutputFormat, "Output format for structured data (yaml or json)",
	).Default("yaml").Allowed("yaml", "json")

	oneshot := &cmd.Command{
		Name:    "oneshot",
		Docline: "Create an ad-hoc prompt from arguments, clipboard or stdin",
		Usage:   "[INSTRUCTIONS...]",
		Aliases: []string{"o"},
		Flags:   qry.BindFlags(),
	}
	oneshot.Runner = func() error {
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
		return run(&qry, func() (config.Prompt, error) { return qry.BuildPrompt(library) })
	}

	force := false
	edit := &cmd.Command{
		Name:    "edit",
		Docline: "Edit files according to instructions",
		Usage:   "[INSTRUCTIONS...] --write FILE...",
		Description: `Only the specified writable files will be modified, any attempt to mutate
anything else will cause an error.

By default only files followed by git and without any unstaged changed are
allowed.
This behavior can be bypassed with --force.`,
		Aliases: []string{"ed"},
		Flags:   qry.BindFlags(),
	}
	edit.Flags.StringSlice("write", &qry.Context.Writable.Files, "Files to edit")
	edit.Flags.StringSlice("write-dirs", &qry.Context.Writable.Dirs, "Dirs to edit (recurses)")
	edit.Flags.Bool("force", &force, "Allow editing files with unstaged changes")
	edit.Runner = func() error {
		return runEdit(&qry, force)
	}

	modelZoo := &cmd.Command{
		Name:    "model-zoo",
		Docline: "Show configured model shortcuts",
		Aliases: []string{"mz"},
		Runner:  printAvailableModels,
	}

	var posargs []string
	details := &cmd.Command{
		Name:    "details",
		Docline: "Print details for a model, or the whole model zoo",
		Usage:   "[MODEL]",
		Flags:   flag.NewParser().Positional(&posargs),
		Aliases: []string{"d"},
	}
	details.Runner = func() error { return printModelDetails(posargs) }

	regen := &cmd.Command{
		Name:        "regen",
		Docline:     "Overwrites the given model zoo YAML with current OpenRouter info",
		Usage:       "OUTPUT-FILE",
		Description: `OUTPUT-FILE will double as the model zoo source if valid.`,
		Flags:       flag.NewParser().Positional(&posargs),
	}
	regen.Runner = func() error {
		return regenModelZoo(&qry, posargs)
	}

	modelZoo.RegisterSubs(details, regen)

	source := ""
	capa := &cmd.Command{
		Name:    "capa",
		Docline: "List a model's OpenRouter capabilities",
		Usage:   "MODEL",
		Flags:   flag.NewParser().Positional(&posargs),
	}
	capa.Flags.String(
		"source", &source, "Capability source: overview, endpointless, endpoints, or rawyaml",
	).Default("overview").Allowed("overview", "endpointless", "endpoints", "rawyaml")
	capa.Runner = func() error { return runCapa(&qry, posargs, source) }

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

//revive:disable:cognitive-complexity,cyclomatic
func runWithHooks(
	qry *config.Query,
	mkprompt func() (config.Prompt, error),
	prepare func() error,
	onResponse func(string) error,
) error {
	if qry.DumpConfig {
		spec := noerr(models.ParseSpec(qry.Model))
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
		dry := prompt.String()
		qry.MaybeCopy(dry)
		return renderAndDisplayMarkdown(dry)
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
	qry.MaybeCopy(response)
	if onResponse != nil {
		return onResponse(response)
	}
	return nil
}

// modelAndAPIKey resolves the model spec and API key for the query's model, applying the reasoning
// override when requested.
func modelAndAPIKey(qry *config.Query) (models.Spec, string, error) {
	spec, err := models.ParseSpec(qry.Model)
	if err != nil {
		return models.Spec{}, "", err
	}
	if err := setModelReasoning(qry, &spec); err != nil {
		return models.Spec{}, "", err
	}
	apiKey, err := qry.APIKeyFor(spec.Provider)
	if err != nil {
		return models.Spec{}, "", err
	}

	return spec, apiKey, nil
}

func setModelReasoning(qry *config.Query, spec *models.Spec) error {
	if qry.Effort == "" {
		return nil
	}

	doset := func() bool {
		reasoning := fun.ZeroDeref(spec.Reasoning)
		if slices.Contains(reasoning.SupportedEfforts, qry.Effort) {
			spec.Reasoning = &models.ReasoningSpec{
				SupportedEfforts: reasoning.SupportedEfforts,
				Effort:           qry.Effort,
			}
			return true
		}
		return false
	}
	if doset() {
		return nil
	}

	// Request effort is apparently not supported, checking upstream to be sure.
	upd, err := openrouter.NewModelUpdater(qry).Updated(*spec)
	if err != nil {
		return err
	}
	*spec = upd

	if doset() {
		return nil
	}

	msg := "reasoning effort %q is not supported by model %q"
	return fmt.Errorf(msg, qry.Effort, spec.QualifiedName())
}
