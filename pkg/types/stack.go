package types

type StackSelection struct {
	Stacks   []ScoredStack `json:"stacks"`
	TopStack string        `json:"top-stack"`
}

type ScoredStack struct {
	Stack
	Score      int      `json:"score"`
	Compatible bool     `json:"compatible"`
	Notes      []string `json:"notes,omitempty"`
}

type Stack struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Vendor      string `yaml:"vendor" json:"vendor"`
	Grade       string `yaml:"grade" json:"grade"`

	Devices   StackDevices `yaml:"devices" json:"devices"`
	Memory    *string      `yaml:"memory" json:"memory"`
	DiskSpace *string      `yaml:"disk-space" json:"disk-space"`

	Components     []string  `yaml:"components" json:"components"`
	Configurations StackConf `yaml:"configurations" json:"configurations"`
}

type StackDevices struct {
	Any []StackDevice `yaml:"any" json:"any"`
	All []StackDevice `yaml:"all" json:"all"`
}

type StackDevice struct {
	Type     string  `yaml:"type" json:"type"`
	VendorId *string `yaml:"vendor-id" json:"vendor-id"`

	// CPUs
	Architectures []string `yaml:"architectures" json:"architectures,omitempty"`
	FamilyIds     []string `yaml:"family-ids" json:"family-ids,omitempty"`
	Flags         []string `yaml:"flags" json:"flags,omitempty"`

	// GPUs
	Bus               *string `yaml:"bus" json:"bus,omitempty"`
	VRam              *string `yaml:"vram" json:"vram,omitempty"`
	ComputeCapability *string `yaml:"compute-capability" json:"compute-capability,omitempty"`
}

type StackConf map[string]interface{}
