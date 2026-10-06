package commands

import (
	"fmt"

	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
	"github.com/canonical/lscompute/pkg/machine/memory"
)

// machineFixture returns a small, hand-built Machine fixture for the named machine.
func machineFixture(name string) (*machine.Machine, error) {
	switch name {
	case "dummy-machine":
		return &machine.Machine{
			CPUs: []cpu.CPU{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"fpu", "vme", "de"},
			}},
			Memory: memory.Memory{TotalRam: 67012501504, TotalSwap: 0},
			Disk: []disk.Disk{{
				Total:     1006451294208,
				Available: 943543738368,
				Path:      "/var/lib/snapd/snaps",
			}},
			PCIDevices: []pci.Device{{
				Bus:                  "pci",
				Slot:                 "0000:00:00.0",
				BusNumber:            0x0,
				DeviceClass:          0x600,
				ProgrammingInterface: new(uint8(0)),
				VendorId:             0x8086,
				DeviceId:             0x4637,
				SubvendorId:          new(uint16(0x103C)),
				SubdeviceId:          new(uint16(0x89C6)),
				FriendlyNames: pci.FriendlyNames{
					VendorName:    "Intel Corporation",
					SubvendorName: "Hewlett-Packard Company",
				}},
			}}, nil

	case "i7-1165G7":
		return &machine.Machine{
			CPUs: []cpu.CPU{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"sse4_2", "f16c", "fma", "avx", "avx2", "avx512f"},
			}},
		}, nil

	case "xps13-7390":
		return &machine.Machine{
			CPUs: []cpu.CPU{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
				Flags:          []string{"sse4_2", "f16c", "fma", "avx", "avx2"},
			}},
			PCIDevices: []pci.Device{{
				Bus:         "pci",
				Slot:        "0000:00:02.0",
				BusNumber:   0x0,
				DeviceClass: 0x300,
				VendorId:    0x8086,
				DeviceId:    0x9B41,
				FriendlyNames: pci.FriendlyNames{
					VendorName: "Intel Corporation",
					DeviceName: "CometLake-U GT2 [UHD Graphics]",
				},
			}},
		}, nil

	case "mustang":
		return &machine.Machine{
			CPUs: []cpu.CPU{{
				Architecture:   "amd64",
				ManufacturerId: "GenuineIntel",
			}},
			Disk: []disk.Disk{{
				Total:     1006451294208,
				Available: 943543738368,
				Path:      "/var/lib/snapd/snaps",
			}},
		}, nil

	case "low-disk-available-machine":
		return &machine.Machine{
			Disk: []disk.Disk{{
				Total:     1006451294208,
				Available: 1024 * 1024 * 1024,
				Path:      "/var/lib/snapd/snaps",
			}},
		}, nil

	case "no-disk-available-machine":
		return &machine.Machine{
			Disk: []disk.Disk{{
				Total:     1006451294208,
				Available: 0,
				Path:      "/var/lib/snapd/snaps",
			}},
		}, nil

	default:
		return nil, fmt.Errorf("no machine fixture for %q", name)
	}
}
