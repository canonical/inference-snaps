package commands

import (
	"fmt"

	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
	"github.com/canonical/lscompute/pkg/machine/memory"
	"github.com/canonical/lscompute/pkg/machine/types"
)

// machineInfoFixture returns a small, hand-built MachineInfo fixture for the named machine.
func machineInfoFixture(name string) (*machine.MachineInfo, error) {
	switch name {
	case "dummy-machine":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"fpu", "vme", "de"},
			}},
			Memory: memory.MemoryInfo{TotalRam: 67012501504, TotalSwap: 0},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1006451294208, Avail: 943543738368},
			},
			Devices: []any{pci.Device{
				Bus:                  "pci",
				Slot:                 "0000:00:00.0",
				BusNumber:            0x0,
				DeviceClass:          0x600,
				ProgrammingInterface: new(uint8(0)),
				VendorId:             0x8086,
				DeviceId:             0x4637,
				SubvendorId:          new(types.HexInt(0x103C)),
				SubdeviceId:          new(types.HexInt(0x89C6)),
				FriendlyNames: pci.FriendlyNames{
					VendorName:    new("Intel Corporation"),
					SubvendorName: new("Hewlett-Packard Company"),
				},
			}},
		}, nil

	case "i7-1165G7":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"sse4_2", "f16c", "fma", "avx", "avx2", "avx512f"},
			}},
		}, nil

	case "xps13-7390":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"sse4_2", "f16c", "fma", "avx", "avx2"},
			}},
			Devices: []any{pci.Device{
				Bus:         "pci",
				Slot:        "0000:00:02.0",
				BusNumber:   0x0,
				DeviceClass: 0x300,
				VendorId:    0x8086,
				DeviceId:    0x9B41,
				FriendlyNames: pci.FriendlyNames{
					VendorName: new("Intel Corporation"),
					DeviceName: new("CometLake-U GT2 [UHD Graphics]"),
				},
			}},
		}, nil

	case "mustang":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
			}},
		}, nil

	default:
		return nil, fmt.Errorf("no machine fixture for %q", name)
	}
}
