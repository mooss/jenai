package main

import (
	"errors"
	"fmt"
	"os"

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

	specup := openrouter.NewModelUpdater(qry)
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

// modelSpecsForRegen reads model specs from the given path, defaulting to the embedded specs if it
// does not exist.
func modelSpecsForRegen(path string) (map[string]models.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		return models.ModelSpecs()
	}

	return models.FromYAML(data)
}
