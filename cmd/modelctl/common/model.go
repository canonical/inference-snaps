package common

import (
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

type ModelDetails struct {
	Name  string `json:"name" yaml:"name"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`

	Description  string   `json:"description" yaml:"description"`
	ModelCardUrl string   `json:"model-card-url" yaml:"model-card-url"`
	Quantization string   `json:"quantization" yaml:"quantization"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	Components []string `json:"components" yaml:"components"`
}

func NewModelDetails(manifest *models.Manifest) (ModelDetails, error) {
	var modelDetails ModelDetails
	modelDetails.Name = manifest.Name
	modelDetails.Alias = manifest.Alias
	modelDetails.Description = manifest.Description
	modelDetails.ModelCardUrl = manifest.ModelCardUrl
	modelDetails.Quantization = manifest.Quantization
	modelDetails.Capabilities = manifest.Capabilities
	modelDetails.Components = manifest.Components

	// Change disk size to largest possible unit representation
	diskSizeBytes, err := utils.StringToBytes(manifest.DiskSize)
	if err != nil {
		return modelDetails, err
	}
	modelDetails.DiskSize = utils.FmtBytesShort(diskSizeBytes)

	return modelDetails, nil
}

func GetModelByNameOrAlias(ctx *Context, modelName string) (*models.Manifest, error) {
	if modelName == "" {
		return nil, fmt.Errorf("model name must not be empty")
	}

	activeEngine, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}
	if activeEngine == "" {
		return nil, ErrNoActiveEngine
	}

	engineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngine)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}

	allModelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}

	// Consider the active engine's models first
	var manifest *models.Manifest
	for _, modelManifest := range allModelManifests {
		if slices.Contains(engineManifest.Model.Options, modelManifest.Name) &&
			(modelManifest.Name == modelName || modelManifest.Alias == modelName) {
			manifest = &modelManifest
			break
		}
	}

	// If the provided name is not one of the active engine's models, warn the user it is not compatible
	if manifest == nil {
		for _, modelManifest := range allModelManifests {
			if modelManifest.Name == modelName || modelManifest.Alias == modelName {
				return nil, fmt.Errorf("model %q is not compatible with the active engine", modelName)
			}
		}
	}

	if manifest == nil {
		return nil, fmt.Errorf("model %q does not exist", modelName)
	}
	return manifest, nil
}

func ModelStatus(ctx *Context) (map[string]string, error) {
	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveModel, err)
	}

	activeModelManifest, err := models.LoadManifest(ctx.ModelsDir, activeModelId)
	if err != nil {
		return nil, fmt.Errorf("loading model manifest: %v", err)
	}

	status := make(map[string]string)
	status["name"] = activeModelManifest.Name

	return status, nil
}
