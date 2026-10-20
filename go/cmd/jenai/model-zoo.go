package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/llm/openrouter"
	"github.com/mooss/jenai/go/models"
	"github.com/mooss/jenai/go/utils"
)

func regenModelZoo(cfg *config.Jenai, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: jenai model-zoo regen OUTPUT-FILE")
	}

	specs, err := models.ModelSpecs()
	if err != nil {
		return err
	}

	specup := specUpdater{
		ctx:    context.Background(),
		apiKey: cfg.APIKeys["openrouter"],
	}
	output := models.Zoo{Models: map[string]models.Spec{}}

	for name, spec := range specs {
		spec, err := specup.update(spec)
		if err != nil {
			return err
		}
		output.Models[name] = spec
	}

	data, err := utils.ToYAML(output)
	if err != nil {
		return fmt.Errorf("model-zoo regen: marshal YAML: %w", err)
	}
	if err := os.WriteFile(args[0], data, 0600); err != nil {
		return fmt.Errorf("model-zoo regen: write %q: %w", args[0], err)
	}
	return nil
}

type specUpdater struct {
	ctx    context.Context
	apiKey string
}

func (su specUpdater) update(spec models.Spec) (models.Spec, error) {
	if spec.Provider != "openrouter" {
		msg := "model-zoo regen: unsupported provider %q for model %q"
		return models.Spec{}, fmt.Errorf(msg, spec.Provider, spec.ShortName)
	}

	modelResponse, err := fetchRecord[openrouter.GetModelResponse](
		su.ctx, modelURL(spec, "model"), su.apiKey, "model",
	)
	if err != nil {
		return models.Spec{}, fmt.Errorf("model-zoo regen: fetch model %q: %w", spec.ShortName, err)
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
	if reasoning := modelResponse.Data.Reasoning; reasoning != nil {
		updated := &models.ReasoningSpec{
			SupportedEfforts: lie.MapFilter(fun.Identity, reasoning.SupportedEfforts),
			DefaultEffort:    fun.ZeroDeref(reasoning.DefaultEffort),
		}

		customDefault := fun.ZeroDeref(spec.Reasoning).DefaultEffort
		if customDefault != "" && slices.Contains(updated.SupportedEfforts, customDefault) {
			updated.DefaultEffort = customDefault // Maintain the custom default.
		}

		if len(updated.SupportedEfforts) > 0 || updated.DefaultEffort != "" {
			res = updated
		}
	}
	return res
}
