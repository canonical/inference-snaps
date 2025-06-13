package types

type HwInfo struct {
	Cpus       []CpuInfo            `json:"cpus,omitempty"`
	Memory     *MemoryInfo          `json:"memory,omitempty"`
	Disk       map[string]*DirStats `json:"disk,omitempty"`
	PciDevices []PciDevice          `json:"pci,omitempty"`
}
