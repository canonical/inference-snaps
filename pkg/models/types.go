package models

import (
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
)

const (
	capabilityText                  string = "text"
	capabilityVision                string = "vision"
	capabilityTools                 string = "tools"
	capabilityThinking              string = "thinking"
	capabilityTextEmbedding         string = "text-embedding"
	capabilityTranscription         string = "transcription"
	capabilityRealtimeTranscription string = "realtime-transcription"
)

const (
	formatGGUF        string = "GGUF"
	formatOpenVINOIR  string = "OpenVINO_IR"
	formatCTranslate2 string = "CTranslate2"
	formatMediaTekDLA string = "MediaTek_DLA"
)

type Manifest struct {
	Name  string `json:"name" yaml:"name"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`

	Description  string   `json:"description" yaml:"description"`
	ModelCardUrl string   `json:"model-card-url" yaml:"model-card-url"`
	Format       string   `json:"format" yaml:"format"`
	Quantization string   `json:"quantization" yaml:"quantization"`
	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	DiskSize string `json:"disk-size" yaml:"disk-size"`

	Components []string `json:"components" yaml:"components"`

	Environment []string `json:"environment" yaml:"environment"`

	Layout map[string]engines.Layout `json:"layout,omitempty" yaml:"layout,omitempty"`
}

func SupportedCapabilities() []string {
	return []string{capabilityText, capabilityVision, capabilityTools, capabilityThinking,
		capabilityTextEmbedding, capabilityTranscription, capabilityRealtimeTranscription}
}

func SupportedFormats() []string {
	return []string{formatGGUF, formatOpenVINOIR, formatCTranslate2, formatMediaTekDLA}
}
