package types

type StackSelection struct {
	Stacks   []ScoredStack `json:"stacks"`
	TopStack string        `json:"top-stack"`
}

type ScoredStack struct {
	Name       string   `json:"name"`
	Score      int      `json:"score"`
	Compatible bool     `json:"compatible"`
	Grade      string   `json:"grade"`
	Notes      []string `json:"notes,omitempty"`
}

type Stack struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Vendor      string `yaml:"vendor"`
	Grade       string `yaml:"grade"`

	Devices   StackDevices `yaml:"devices"`
	Memory    *string      `yaml:"memory"`
	DiskSpace *string      `yaml:"disk-space"`

	Components     []string  `yaml:"components"`
	Configurations StackConf `yaml:"configurations"`
}

type StackDevices struct {
	Any []StackDevice `yaml:"any"`
	All []StackDevice `yaml:"all"`
}

type StackDevice struct {
	Type     string  `yaml:"type"`
	VendorId *string `yaml:"vendor-id"`

	// CPUs
	Architectures []string `yaml:"architectures"`
	FamilyIds     []string `yaml:"family-ids"`
	Flags         []string `yaml:"flags"`

	// GPUs
	Bus               *string `yaml:"bus"`
	VRam              *string `yaml:"vram"`
	ComputeCapability *string `yaml:"compute-capability"`
}

type StackConf map[string]interface{}
