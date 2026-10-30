package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/bagend/go/fun/eager/mae"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/utils"
)

func printModelDetails(args []string) error {
	if len(args) > 1 {
		return errors.New("usage: jenai model-zoo details [MODEL]")
	}

	if len(args) == 0 {
		specs, err := models.ModelSpecs()
		if err != nil {
			return err
		}
		prettyPrint(models.Zoo{Models: specs})
		return nil
	}

	spec, err := models.Get(args[0])
	if err != nil {
		return err
	}
	prettyPrint(spec)
	return nil
}

func regenModelZoo(qry *config.Query, args []string) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = fmt.Errorf("model-zoo regen: %w", retErr)
		}
	}()

	if len(args) != 1 {
		return errors.New("usage: jenai model-zoo regen OUTPUT-FILE")
	}

	specs, err := modelSpecsForRegen(args[0])
	if err != nil {
		return err
	}

	specup := NewSpecUpdater(qry)
	updated, err := mae.MapVE(specup.Updated, specs)
	if err != nil {
		return err
	}
	output := models.Zoo{Models: updated}

	data, err := utils.ToYAML(output)
	if err != nil {
		return fmt.Errorf("marshal YAML: %w", err)
	}
	if err := os.WriteFile(args[0], data, 0600); err != nil {
		return fmt.Errorf("write %q: %w", args[0], err)
	}
	return nil
}

func modelSpecsForRegen(path string) (map[string]models.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		return models.ModelSpecs()
	}

	specs, err := models.FromYAML(data)
	if err == nil {
		return specs, nil
	}

	warn := "could not load model specs from %q: %v; using embedded specs\n"
	fmt.Fprintf(os.Stderr, warn, path, err)
	return models.ModelSpecs()
}

///////////////////
// Spec updating //

type SpecUpdater struct {
	ctx    context.Context
	apiKey string
}

func NewSpecUpdater(qry *config.Query) SpecUpdater {
	return SpecUpdater{
		ctx:    context.Background(),
		apiKey: qry.APIKeys["openrouter"],
	}
}

func (su SpecUpdater) Updated(spec models.Spec) (models.Spec, error) {
	if spec.Provider != "openrouter" {
		msg := "unsupported provider %q for model %q"
		return models.Spec{}, fmt.Errorf(msg, spec.Provider, spec.ShortName)
	}

	modelResponse, err := openrouter.FetchModel(su.ctx, spec, su.apiKey)
	if err != nil {
		return models.Spec{}, fmt.Errorf("fetch model %q: %w", spec.ShortName, err)
	}

	return models.Spec{
		Provider:  spec.Provider,
		Author:    spec.Author,
		Model:     spec.Model,
		Reasoning: buildReasoningSpec(spec, modelResponse),
	}, nil
}

func buildReasoningSpec(
	spec models.Spec, modelResponse openrouter.GetModelResponse,
) *models.ReasoningSpec {
	var res *models.ReasoningSpec
	reasoning := modelResponse.Data.Reasoning
	if reasoning == nil {
		return res
	}

	updated := &models.ReasoningSpec{Effort: fun.ZeroDeref(reasoning.DefaultEffort)}
	if reasoning.SupportedEfforts == nil {
		updated.SupportedEfforts = strings.Split("max xhigh high medium low minimal", " ")
		if !reasoning.Mandatory {
			updated.SupportedEfforts = append(updated.SupportedEfforts, "none")
		}
	} else {
		updated.SupportedEfforts = lie.MapFilter(fun.Identity, *reasoning.SupportedEfforts)
	}

	customDefault := fun.ZeroDeref(spec.Reasoning).Effort
	if customDefault != "" && slices.Contains(updated.SupportedEfforts, customDefault) {
		updated.Effort = customDefault // Maintain the custom default.
	}

	if len(updated.SupportedEfforts) > 0 || updated.Effort != "" {
		res = updated
	}

	return res
}
