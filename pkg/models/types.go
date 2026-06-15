package models

import "github.com/canonical/inference-snaps-cli/v2/pkg/types"

type Manifest struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`

	Description  string   `json:"description" yaml:"description"`
	ModelCardUrl string   `json:"model-card-url" yaml:"model-card-url"`
	Quantization string   `json:"quantization" yaml:"quantization"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	Components []string `json:"components" yaml:"components"`

	Environment []string `json:"environment" yaml:"environment"`

	Layout map[string]types.Layout `json:"layout,omitempty" yaml:"layout,omitempty"`
}
