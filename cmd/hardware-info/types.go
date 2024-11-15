package main

import (
	"github.com/canonical/hardware-info/cpu"
	"github.com/canonical/hardware-info/disk"
	"github.com/canonical/hardware-info/memory"
)

type HwInfo struct {
	Cpu    *cpu.CpuInfo              `json:"cpu,omitempty"`
	Memory *memory.MemoryInfo        `json:"memory,omitempty"`
	Disk   map[string]*disk.DirStats `json:"disk,omitempty"`
}
