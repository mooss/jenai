package openrouter

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/mooss/bagend/go/fun"
	"github.com/mooss/bagend/go/fun/eager/lie"
	"github.com/mooss/jenai/go/config"
	"github.com/mooss/jenai/go/models"
)

type ModelUpdater struct {
	ctx    context.Context
	apiKey string
}

func NewModelUpdater(qry *config.Query) ModelUpdater {
	return ModelUpdater{
		ctx:    context.Background(),
		apiKey: qry.APIKeys["openrouter"],
	}
}

func (su ModelUpdater) Updated(spec models.Spec) (models.Spec, error) {
	if spec.Provider != "openrouter" {
		msg := "unsupported provider %q for model %q"
		return models.Spec{}, fmt.Errorf(msg, spec.Provider, spec.ShortName)
	}

	modelResponse, err := FetchModel(su.ctx, spec, su.apiKey)
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
	spec models.Spec, modelResponse GetModelResponse,
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
