package types

type HwInfo struct {
	Cpus       []CpuInfo           `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	Memory     MemoryInfo          `json:"memory,omitempty" yaml:"memory,omitempty"`
	Disk       map[string]DirStats `json:"disk,omitempty" yaml:"disk,omitempty"`
	PciDevices []PciDevice         `json:"pci,omitempty" yaml:"pci"`
}
