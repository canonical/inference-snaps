package types

import (
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/cpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/disk"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/gpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/memory"
)

type HwInfo struct {
	Cpu    *cpu.CpuInfo              `json:"cpu,omitempty"`
	Memory *memory.MemoryInfo        `json:"memory,omitempty"`
	Disk   map[string]*disk.DirStats `json:"disk,omitempty"`
	Gpus   []gpu.Gpu                 `json:"gpu,omitempty"`
}
