package engines

import "github.com/canonical/famous-models-cli/pkg/types"

type ScoredManifest struct {
	Manifest   `yaml:",inline"`
	Score      int      `yaml:"score" json:"score"`
	Compatible bool     `yaml:"compatible" json:"compatible"`
	Notes      []string `yaml:"notes,omitempty" json:"notes,omitempty"`
}

type Manifest struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Vendor      string `yaml:"vendor" json:"vendor"`
	Grade       string `yaml:"grade" json:"grade"`

	Devices   Devices `yaml:"devices" json:"devices"`
	Memory    *string `yaml:"memory,omitempty" json:"memory"`
	DiskSpace *string `yaml:"disk-space,omitempty" json:"disk-space"`

	Components     []string       `yaml:"components" json:"components"`
	Configurations Configurations `yaml:"configurations" json:"configurations"`
}

type Devices struct {
	Any []Device `yaml:"any,omitempty" json:"any"`
	All []Device `yaml:"all,omitempty" json:"all"`
}

type Device struct {
	Type string `yaml:"type,omitempty" json:"type,omitempty"` // cpu, gpu, npu or nil
	Bus  string `yaml:"bus,omitempty" json:"bus,omitempty"`   // pci, usb or nil

	// CPUs
	Architecture *string `yaml:"architecture,omitempty" json:"architecture,omitempty"`

	// CPU x86
	ManufacturerId *string  `yaml:"manufacturer-id,omitempty" json:"manufacturer-id,omitempty"`
	Flags          []string `yaml:"flags,omitempty" json:"flags,omitempty"`

	// CPU arm64
	ImplementerId *types.HexInt `yaml:"implementer-id,omitempty" json:"implementer-id,omitempty"`
	PartNumber    *types.HexInt `yaml:"part-number,omitempty" json:"part-number,omitempty"`
	Features      []string      `yaml:"features,omitempty" json:"features,omitempty"`

	// PCI
	VendorId *types.HexInt `yaml:"vendor-id,omitempty" json:"vendor-id,omitempty"`
	DeviceId *types.HexInt `yaml:"device-id,omitempty" json:"device-id,omitempty"`

	// GPU additional properties
	VRam              *string `yaml:"vram,omitempty" json:"vram,omitempty"`
	ComputeCapability *string `yaml:"compute-capability,omitempty" json:"compute-capability,omitempty"`

	// NPU
	// no additional properties for now
}

type Configurations map[string]interface{}
