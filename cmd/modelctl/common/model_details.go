package common

import (
	"github.com/canonical/inference-snaps-cli/v2/pkg/models"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

type ModelDetails struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`

	Description  string   `json:"description" yaml:"description"`
	ModelCardUrl string   `json:"model-card-url" yaml:"model-card-url"`
	Quantization string   `json:"quantization" yaml:"quantization"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	Components []string `json:"components" yaml:"components"`
}

func NewModelDetails(manifest *models.Manifest) (ModelDetails, error) {
	var modelDetails ModelDetails
	modelDetails.ID = manifest.ID
	modelDetails.Name = manifest.Name
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
