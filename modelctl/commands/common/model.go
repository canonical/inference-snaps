package common

import (
	"fmt"
	"slices"

	"github.com/canonical/inference-snaps-cli/v2/internal/constants"
	"github.com/canonical/inference-snaps-cli/v2/internal/engines"
	"github.com/canonical/inference-snaps-cli/v2/internal/models"
	selectorpci "github.com/canonical/inference-snaps-cli/v2/internal/selector/pci"
	"github.com/canonical/inference-snaps-cli/v2/internal/utils"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
)

type ModelDetails struct {
	Name  string `json:"name" yaml:"name"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`

	Description  string   `json:"description,omitempty" yaml:"description,omitempty"`
	ModelCardUrl string   `json:"model-card-url,omitempty" yaml:"model-card-url,omitempty"`
	Format       string   `json:"format,omitempty" yaml:"format,omitempty"`
	Quantization string   `json:"quantization,omitempty" yaml:"quantization,omitempty"`
	Capabilities []string `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`

	DiskSize       string `json:"disk-size" yaml:"disk-size"`
	RequiredMemory string `json:"required-memory,omitempty" yaml:"required-memory,omitempty"`

	Components []string `json:"components,omitempty" yaml:"components,omitempty"`

	CompatibleEngines   []string `json:"compatible-engines,omitempty" yaml:"compatible-engines,omitempty"`
	CompatibilityIssues []string `json:"compatibility-issues,omitempty" yaml:"compatibility-issues,omitempty"`
}

func (m *ModelDetails) fillIncompatibilityIssues(report models.CompatibilityReport) {
	var issues []string
	if report.TotalVram < report.RequiredVram {
		issues = append(issues, "insufficient VRAM")
	} else if report.TotalMemory < report.RequiredMemory {
		issues = append(issues, "insufficient system memory")
	}

	if !report.CompatibleDisk {
		issues = append(issues, "insufficient disk space")
	}

	m.CompatibilityIssues = issues
}

func NewModelDetails(manifest *models.Manifest) (ModelDetails, error) {
	var modelDetails ModelDetails
	modelDetails.Name = manifest.Name
	modelDetails.Alias = manifest.Alias
	modelDetails.Description = manifest.Description
	modelDetails.ModelCardUrl = manifest.ModelCardUrl
	modelDetails.Format = manifest.Format
	modelDetails.Quantization = manifest.Quantization
	modelDetails.Capabilities = manifest.Capabilities
	modelDetails.Components = manifest.Components

	// Change disk size to largest possible unit representation
	diskSize, err := utils.NormalizeFormattedBytes(manifest.DiskSize)
	if err != nil {
		return modelDetails, ErrInsufficientDiskSpaceForModel
	}
	modelDetails.DiskSize = diskSize

	if manifest.RequiredMemory != "" {
		requiredMemory, err := utils.NormalizeFormattedBytes(manifest.RequiredMemory)
		if err != nil {
			return modelDetails, err
		}
		modelDetails.RequiredMemory = requiredMemory
	}

	return modelDetails, nil
}

func GetModelManifestByNameOrAlias(ctx *Context, modelName string) (*models.Manifest, error) {
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

func GetModelDetailsByNameOrAlias(ctx *Context, modelName string, machine *machine.Machine) (*ModelDetails, error) {
	modelManifest, err := GetModelManifestByNameOrAlias(ctx, modelName)
	if err != nil {
		return nil, err
	}
	modelDetails, err := NewModelDetails(modelManifest)
	if err != nil {
		return nil, err
	}
	compatibleEngines, err := GetCompatibleEnginesByModelName(ctx, modelDetails.Name)
	if err != nil {
		return nil, err
	}
	modelDetails.CompatibleEngines = compatibleEngines
	activeEngine, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}
	if activeEngine == "" {
		return nil, ErrNoActiveEngine
	}
	activeEngineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngine)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}
	modelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}
	manifestsByName := make(map[string]models.Manifest, len(modelManifests))
	for _, manifest := range modelManifests {
		manifestsByName[manifest.Name] = manifest
	}
	scoredModels, err := ScoreModels(activeEngineManifest.Model.Options, manifestsByName, machine, *activeEngineManifest)
	if err != nil {
		return nil, err
	}
	for _, model := range scoredModels {
		if model.Name == modelDetails.Name {
			modelDetails.fillIncompatibilityIssues(model.CompatibilityReport)
			return &modelDetails, nil
		}
	}
	return nil, fmt.Errorf("model %q not found in scored models", modelDetails.Name)
}

func GetCompatibleEnginesByModelName(ctx *Context, modelName string) ([]string, error) {
	allEngineManifests, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}

	compatibleEngines := []string{}
	for _, engineManifest := range allEngineManifests {
		if slices.Contains(engineManifest.Model.Options, modelName) {
			compatibleEngines = append(compatibleEngines, engineManifest.Name)
		}
	}
	return compatibleEngines, nil
}

func ModelStatus(ctx *Context) (map[string]string, error) {
	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveModel, err)
	}

	if activeModelId == "" {
		return nil, ErrNoActiveModel
	}

	activeModelManifest, err := models.LoadManifest(ctx.ModelsDir, activeModelId)
	if err != nil {
		return nil, fmt.Errorf("loading model manifest: %v", err)
	}

	status := make(map[string]string)
	status["name"] = activeModelManifest.Name

	return status, nil
}

func GetAllModels(ctx *Context) ([]ModelDetails, error) {
	allModelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}

	allEngineManifests, err := engines.LoadManifests(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LoadingEngineManifest, err)
	}

	allModelsWithEngines := []ModelDetails{}
	for _, modelManifest := range allModelManifests {
		outputModel, err := NewModelDetails(&modelManifest)
		if err != nil {
			return nil, fmt.Errorf("creating model details for model %s: %v", modelManifest.Name, err)
		}
		compatibleEngines := []string{}
		for _, engineManifest := range allEngineManifests {
			if slices.Contains(engineManifest.Model.Options, modelManifest.Name) {
				compatibleEngines = append(compatibleEngines, engineManifest.Name)
			}
		}
		outputModel.CompatibleEngines = compatibleEngines
		allModelsWithEngines = append(allModelsWithEngines, outputModel)
	}

	return allModelsWithEngines, nil
}

func ScoreModels(modelOptions []string, manifests map[string]models.Manifest, machine *machine.Machine, activeEngineManifest engines.Manifest) ([]models.ScoredManifest, error) {
	availableDiskSpace, err := availableDiskSpace(machine)
	if err != nil {
		return nil, err
	}

	systemMemory, err := getSystemMemory(machine)
	if err != nil {
		return nil, err
	}

	enginesRequiredGpus := GetEngineGpuDevices(activeEngineManifest)
	vRam, isUnifiedMemory, err := getTopGpuVram(machine, enginesRequiredGpus)
	if err != nil {
		return nil, err
	}

	scoredModels := make([]models.ScoredManifest, 0, len(modelOptions))
	for _, modelID := range modelOptions {
		manifest, ok := manifests[modelID]
		if !ok {
			return nil, fmt.Errorf("model manifest not found: %s", modelID)
		}

		diskSize, err := utils.StringToBytes(manifest.DiskSize)
		if err != nil {
			return nil, fmt.Errorf("parsing disk size for model %q: %w", modelID, err)
		}

		var modelRequiredMemory uint64
		if manifest.RequiredMemory == "" {
			modelRequiredMemory = diskSize
		} else {
			modelRequiredMemory, err = utils.StringToBytes(manifest.RequiredMemory)
			if err != nil {
				return nil, fmt.Errorf("parsing required memory for model %q: %w", modelID, err)
			}
		}

		requiredVram := uint64(0)
		requiredMemory := uint64(0)
		if len(enginesRequiredGpus) > 0 && !isUnifiedMemory {
			requiredVram = modelRequiredMemory
		} else {
			requiredMemory = modelRequiredMemory
		}

		report := models.CompatibilityReport{
			CompatibleDisk:     diskSize <= availableDiskSpace,
			RequiredDiskSpace:  diskSize,
			AvailableDiskSpace: availableDiskSpace,
			RequiredMemory:     requiredMemory,
			TotalMemory:        systemMemory,
			RequiredVram:       requiredVram,
			TotalVram:          vRam,
		}

		score := uint64(0)
		if report.Compatible() {
			score = diskSize
		}

		scoredModels = append(scoredModels, models.ScoredManifest{
			Manifest:            manifest,
			Score:               score,
			CompatibilityReport: report,
		})
	}

	return scoredModels, nil
}

// SelectModel picks a model that fits the available disk and memory: the preferred one if it fits, otherwise the largest fitting model.
// Returns ErrInsufficientDiskSpaceForModel or ErrInsufficientMemoryForModel when none fit.
func SelectModel(ctx *Context, engineManifest engines.Manifest, preferredModel string, machine *machine.Machine) (string, []models.ScoredManifest, error) {
	modelManifests, err := models.LoadManifests(ctx.ModelsDir)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", LoadingModelManifests, err)
	}
	manifestsByName := make(map[string]models.Manifest, len(modelManifests))
	for _, manifest := range modelManifests {
		manifestsByName[manifest.Name] = manifest
	}

	scoredModels, err := ScoreModels(engineManifest.Model.Options, manifestsByName, machine, engineManifest)
	if err != nil {
		return "", nil, err
	}

	if preferredModel != "" {
		for _, model := range scoredModels {
			if model.Name == preferredModel && model.CompatibilityReport.Compatible() {
				return preferredModel, scoredModels, nil
			}
		}
	}

	selected := ""
	selectedSize := uint64(0)
	for _, model := range scoredModels {
		size := model.CompatibilityReport.RequiredDiskSpace
		if model.CompatibilityReport.Compatible() && (selected == "" || size > selectedSize) {
			selected = model.Name
			selectedSize = size
		}
	}

	if selected != "" {
		return selected, scoredModels, nil
	}

	for _, model := range scoredModels {
		if model.CompatibilityReport.CompatibleDisk {
			return "", scoredModels, ErrInsufficientMemoryForModel
		}
	}
	return "", scoredModels, ErrInsufficientDiskSpaceForModel
}

func availableDiskSpace(machine *machine.Machine) (uint64, error) {
	for _, disk := range machine.Disk {
		if disk.Path == constants.SnapStoragePath {
			return disk.Available, nil
		}
	}
	return 0, fmt.Errorf("disk information unavailable for %s", constants.SnapStoragePath)
}

func getTopGpuVram(machine *machine.Machine, gpuDevices []engines.Device) (uint64, bool, error) {
	isUnifiedMemory := false

	var bestDevice *pci.Device
	bestScore := 0
	for _, gpuDevice := range gpuDevices {
		hostDevice, score := selectorpci.BestMatch(gpuDevice, machine)
		if hostDevice != nil && score > bestScore {
			bestDevice = hostDevice
			bestScore = score
		}
	}
	if bestDevice == nil {
		return 0, isUnifiedMemory, nil
	}
	vram, ok := bestDevice.AdditionalProperties["vram"]
	if !ok {
		return 0, isUnifiedMemory, nil
	} else if vram == nil {
		// assume unified memory
		isUnifiedMemory = true
		return 0, isUnifiedMemory, nil
	}
	vramValue, ok := vram.(uint64)
	if !ok {
		return 0, isUnifiedMemory, fmt.Errorf("parsing vram %q for pci device %s", vram, bestDevice.Slot)
	}
	return vramValue, isUnifiedMemory, nil
}

func getSystemMemory(machine *machine.Machine) (uint64, error) {
	if machine.Memory.TotalRam == 0 {
		return 0, ErrSystemMemoryUnavailable
	}
	return machine.Memory.TotalRam, nil
}
