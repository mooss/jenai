// Package models defines models that jenai can use and utilities around them.
package models

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/mooss/jenai/go/utils"
)

// EmbeddedBytes holds the content of the embedded models.yaml file.
//
//go:embed models.yaml
var EmbeddedBytes []byte

// Spec describes a model: where to find it (provider, author, model ID) and optional reasoning
// parameters.
type Spec struct {
	// ShortName of the model, e.g., "r1".
	ShortName string `yaml:"short_name,omitempty"`

	// Provider of the model, e.g., "openrouter".
	Provider string `yaml:"provider"`

	// Author of the model, e.g., "deepseek".
	Author string `yaml:"author"`

	// Model identifier, e.g., "deepseek-r1-0528:nitro".
	Model string `yaml:"model"`

	// Reasoning parameters.
	Reasoning *ReasoningSpec `yaml:"reasoning,omitempty"`
}

// ParseSpec creates a Spec from a model name.
// The name can be a short name (looked up in the model zoo) or a qualified name in
// provider:author/model format.
func ParseSpec(name string) (Spec, error) {
	// Qualified name in provider:author/model format.
	provider, rest, found := strings.Cut(name, ":")
	if !found {
		return Get(name)
	}

	author, model, found := strings.Cut(rest, "/")
	if !found {
		return Get(name)
	}

	return Spec{Provider: provider, Author: author, Model: model}, nil
}

// ReasoningSpec stores optional reasoning settings for a predefined model.
type ReasoningSpec struct {
	SupportedEfforts []string `yaml:"supported_efforts,omitempty"`
	Effort           string   `yaml:"effort,omitempty"`
}

// Zoo is the top-level structure of the models YAML file.
type Zoo struct {
	Models map[string]Spec `yaml:"models"`
}

// loadModels lazily parses the embedded models.yaml file exactly once.
var loadModels = utils.OnceErr(func() (map[string]Spec, error) { return FromYAML(EmbeddedBytes) })

// FromYAML parses model specifications from YAML data, keyed by short name.
func FromYAML(data []byte) (map[string]Spec, error) {
	res, err := utils.FromYAML[Zoo](data)
	if err != nil {
		return nil, fmt.Errorf("failed to load model zoo from YAML: %w", err)
	}

	for shortName, spec := range res.Models {
		spec.ShortName = shortName
		res.Models[shortName] = spec
	}

	return res.Models, nil
}

// ModelSpecs returns the map of model short names to their specifications.
func ModelSpecs() (map[string]Spec, error) {
	return loadModels()
}

// Get returns the specification for the given short name.
func Get(shortName string) (Spec, error) {
	specs, err := ModelSpecs()
	if err != nil {
		return Spec{}, err
	}
	spec, ok := specs[shortName]
	if !ok {
		return Spec{}, fmt.Errorf("unknown model: %s", shortName)
	}
	return spec, nil
}

// ID returns the model identifier as provider APIs expect it (author/model).
func (sp Spec) ID() string {
	return sp.Author + "/" + sp.Model
}

// QualifiedName returns the full model name formatted as provider:author/model.
func (sp Spec) QualifiedName() string {
	return fmt.Sprintf("%s:%s/%s", sp.Provider, sp.Author, sp.Model)
}
