package common

type StackResult struct {
	Name       string   `json:"name"`
	Components []string `json:"components"`
	Score      float64  `json:"score"`
}

type Stack struct {
	Name           string                 `yaml:"name"`
	Description    string                 `yaml:"description"`
	Maintainer     string                 `yaml:"maintainer"`
	Devices        StackDevices           `yaml:"devices"`
	Memory         *string                `yaml:"memory"`
	DiskSpace      *string                `yaml:"disk-space"`
	Components     []string               `yaml:"components"`
	Configurations map[string]interface{} `yaml:"configurations"`
}

type StackDevices struct {
	Any []StackDevice `yaml:"any"`
	All []StackDevice `yaml:"all"`
}

type StackDevice struct {
	Type     string   `yaml:"type"`
	Bus      *string  `yaml:"bus"`
	VendorId *string  `yaml:"vendor-id"`
	VRam     *string  `yaml:"vram"`
	Flags    []string `yaml:"flags"`
}
