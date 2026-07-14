package models

import (
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
)

type Manifest struct {
	Name  string `json:"name" yaml:"name"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`

	Description  string   `json:"description" yaml:"description"`
	ModelCardUrl string   `json:"model-card-url" yaml:"model-card-url"`
	Quantization string   `json:"quantization" yaml:"quantization"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	Components []string `json:"components" yaml:"components"`

	Environment []string `json:"environment" yaml:"environment"`

	Layout map[string]engines.Layout `json:"layout,omitempty" yaml:"layout,omitempty"`
}
