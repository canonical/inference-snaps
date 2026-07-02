package selector

import (
	"fmt"

	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/cpu"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
	"github.com/canonical/lscompute/pkg/machine/disk"
	"github.com/canonical/lscompute/pkg/machine/memory"
	"github.com/canonical/lscompute/pkg/machine/types"
)

// machineInfoFixture returns an in-code MachineInfo fixture for the named machine.
// PCI devices are trimmed to those relevant to engine matching (GPUs,
// accelerators, or devices carrying additional properties).
func machineInfoFixture(name string) (*machine.MachineInfo, error) {
	switch name {
	case "ampere-one-m-banshee-12":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0xC0,
					PartNumber:    0xAC4,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "jscvt", "fcma", "lrcpc", "dcpop", "sha3", "sm3", "sm4", "asimddp", "sha512", "asimdfhm", "dit", "uscat", "ilrcpc", "flagm", "ssbs", "sb", "paca", "pacg", "dcpodp", "flagm2", "frint", "i8mm", "bf16", "rng", "bti", "ecv"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0004:02:00.0",
					BusNumber:            0x2,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x1A03,
					DeviceId:             0x2000,
					SubvendorId:          new(types.HexInt(0x1A03)),
					SubdeviceId:          new(types.HexInt(0x2000)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("ASPEED Technology, Inc."),
						DeviceName:    new("ASPEED Graphics Family"),
						SubvendorName: new("ASPEED Technology, Inc."),
						SubdeviceName: new("Onboard Graphics"),
					},
				},
			},
		}, nil
	case "ampere-one-siryn":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0xC0,
					PartNumber:    0xAC3,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "jscvt", "fcma", "lrcpc", "dcpop", "sha3", "asimddp", "sha512", "asimdfhm", "dit", "uscat", "ilrcpc", "flagm", "ssbs", "sb", "paca", "pacg", "dcpodp", "flagm2", "frint", "i8mm", "bf16", "rng", "bti", "ecv"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
		}, nil
	case "ampere-one-x-banshee-8":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0xC0,
					PartNumber:    0xAC4,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "jscvt", "fcma", "lrcpc", "dcpop", "sha3", "sm3", "sm4", "asimddp", "sha512", "asimdfhm", "dit", "uscat", "ilrcpc", "flagm", "ssbs", "sb", "paca", "pacg", "dcpodp", "flagm2", "frint", "i8mm", "bf16", "rng", "bti", "ecv"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0004:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x1A03,
					DeviceId:             0x2000,
					SubvendorId:          new(types.HexInt(0x1A03)),
					SubdeviceId:          new(types.HexInt(0x2000)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("ASPEED Technology, Inc."),
						DeviceName:    new("ASPEED Graphics Family"),
						SubvendorName: new("ASPEED Technology, Inc."),
						SubdeviceName: new("Onboard Graphics"),
					},
				},
			},
		}, nil
	case "asus-ux301l":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "cpuid_fault", "epb", "pti", "ssbd", "ibrs", "ibpb", "stibp", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "xsaveopt", "dtherm", "ida", "arat", "pln", "pts", "vnmi", "md_clear", "flush_l1d"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 7663951872, TotalSwap: 4294963200},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 124313583616, Avail: 61499084800},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0xA2E,
					SubvendorId:          new(types.HexInt(0x1043)),
					SubdeviceId:          new(types.HexInt(0x13BD)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Haswell-ULT Integrated Graphics Controller"),
						SubvendorName: new("ASUSTeK Computer Inc."),
					},
				},
			},
		}, nil
	case "hp-pavilion-15-cs-3037nl":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "avx512f", "avx512dq", "rdseed", "adx", "smap", "avx512ifma", "clflushopt", "intel_pt", "avx512cd", "sha_ni", "avx512bw", "avx512vl", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "vnmi", "avx512vbmi", "umip", "pku", "ospke", "avx512_vbmi2", "gfni", "vaes", "vpclmulqdq", "avx512_vnni", "avx512_bitalg", "avx512_vpopcntdq", "rdpid", "fsrm", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 16537874432, TotalSwap: 19327348736},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 498347700224, Avail: 61093965824},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x8A52,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x86E2)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Iris Plus Graphics G7"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:06:00.0",
					BusNumber:            0x6,
					DeviceClass:          0x302,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x10DE,
					DeviceId:             0x1D13,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x86E2)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("NVIDIA Corporation"),
						DeviceName:    new("GP108M [GeForce MX250]"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
					AdditionalProperties: map[string]string{
						"computeCapability": "6.1",
						"vram":              "2147483648",
					},
				},
			},
		}, nil
	case "hp-proliant-rl300-gen11-altra":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0C,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp", "ssbs"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 33134252032, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 470888099840, Avail: 435367202816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0006:01:00.1",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x102B,
					DeviceId:             0x538,
					SubvendorId:          new(types.HexInt(0x1590)),
					SubdeviceId:          new(types.HexInt(0xE4)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Matrox Electronics Systems Ltd."),
						DeviceName:    new("MGA G200eH3"),
						SubvendorName: new("Hewlett Packard Enterprise"),
						SubdeviceName: new("iLO5 VGA"),
					},
				},
				pci.Device{
					Bus:         "pci",
					Slot:        "0000:00:00.0",
					BusNumber:   0x0,
					DeviceClass: 0x600,
					VendorId:    0x1DEF,
					DeviceId:    0xE100,
					FriendlyNames: pci.FriendlyNames{
						VendorName: new("Ampere Computing, LLC"),
						DeviceName: new("Altra PCI Express Root Complex A"),
					},
				},
			},
		}, nil
	case "hp-proliant-rl300-gen11-altra-max":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0C,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp", "ssbs"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 66839527424, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 470888099840, Avail: 425825845248},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0006:01:00.1",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x102B,
					DeviceId:             0x538,
					SubvendorId:          new(types.HexInt(0x1590)),
					SubdeviceId:          new(types.HexInt(0xE4)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Matrox Electronics Systems Ltd."),
						DeviceName:    new("MGA G200eH3"),
						SubvendorName: new("Hewlett Packard Enterprise"),
						SubdeviceName: new("iLO5 VGA"),
					},
				},
				pci.Device{
					Bus:         "pci",
					Slot:        "0000:00:00.0",
					BusNumber:   0x0,
					DeviceClass: 0x600,
					VendorId:    0x1DEF,
					DeviceId:    0xE100,
					FriendlyNames: pci.FriendlyNames{
						VendorName: new("Ampere Computing, LLC"),
						DeviceName: new("Altra PCI Express Root Complex A"),
					},
				},
			},
		}, nil
	case "hp-zbook-power-16-inch-g11":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "intel_ppin", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 33047601152, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1003736440832, Avail: 876099534848},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x7DD5,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x8CA0)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Meteor Lake-P [Intel Graphics]"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
					AdditionalProperties: map[string]string{
						"vram": "30592884736",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:0b.0",
					BusNumber:            0x0,
					DeviceClass:          0x1200,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x7D1D,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x8CA0)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Meteor Lake NPU"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0x302,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x10DE,
					DeviceId:             0x28B9,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x8CA3)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("NVIDIA Corporation"),
						DeviceName:    new("AD107GLM [RTX 1000 Ada Generation Laptop GPU]"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
				},
			},
		}, nil
	case "hp-zbook-i712850HX+RadeonPROW6600M":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "tme", "rdpid", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities", "ibpb_exit_to_user"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 67012501504, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1006451294208, Avail: 943543738368},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x4688,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x89C6)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Alder Lake-HX GT1 [UHD Graphics 770]"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x1002,
					DeviceId:             0x73E1,
					SubvendorId:          new(types.HexInt(0x103C)),
					SubdeviceId:          new(types.HexInt(0x89C6)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Advanced Micro Devices, Inc. [AMD/ATI]"),
						DeviceName:    new("Navi 23 WKS-XM [Radeon PRO W6600M]"),
						SubvendorName: new("Hewlett-Packard Company"),
					},
					AdditionalProperties: map[string]string{
						"microarchitecture": "gfx1032",
						"vram":              "8573157376",
					},
				},
			},
		}, nil
	case "i5-3570k+arc-a580+gtx1080ti":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "rdtscp", "lm", "constant_tsc", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "est", "tm2", "ssse3", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "cpuid_fault", "epb", "pti", "ssbd", "ibrs", "ibpb", "stibp", "fsgsbase", "smep", "erms", "xsaveopt", "dtherm", "ida", "arat", "pln", "pts", "md_clear", "flush_l1d"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 16652423168, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1005867986944, Avail: 862075604992},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x162,
					SubvendorId:          new(types.HexInt(0x1458)),
					SubdeviceId:          new(types.HexInt(0xD000)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("IvyBridge GT2 [HD Graphics 4000]"),
						SubvendorName: new("Gigabyte Technology Co., Ltd"),
					},
					AdditionalProperties: map[string]string{
						"vram": "17179869184",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x10DE,
					DeviceId:             0x1B06,
					SubvendorId:          new(types.HexInt(0x10DE)),
					SubdeviceId:          new(types.HexInt(0x1B06)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("NVIDIA Corporation"),
						DeviceName:    new("GP102 [GeForce GTX 1080 Ti]"),
						SubvendorName: new("NVIDIA Corporation"),
					},
					AdditionalProperties: map[string]string{
						"compute_capability": "6.1",
						"vram":               "11811160064",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:04:00.0",
					BusNumber:            0x4,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x56A2,
					SubvendorId:          new(types.HexInt(0x1849)),
					SubdeviceId:          new(types.HexInt(0x6003)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("DG2 [Arc A580]"),
						SubvendorName: new("ASRock Incorporation"),
					},
					AdditionalProperties: map[string]string{
						"vram": "8096681984",
					},
				},
			},
		}, nil
	case "i7-1165G7":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "cat_l2", "invpcid_single", "cdp_l2", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "vnmi", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdt_a", "avx512f", "avx512dq", "rdseed", "adx", "smap", "avx512ifma", "clflushopt", "clwb", "intel_pt", "avx512cd", "sha_ni", "avx512bw", "avx512vl", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "avx512vbmi", "umip", "pku", "ospke", "avx512_vbmi2", "gfni", "vaes", "vpclmulqdq", "avx512_vnni", "avx512_bitalg", "avx512_vpopcntdq", "rdpid", "movdiri", "movdir64b", "fsrm", "avx512_vp2intersect", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 16493559808, TotalSwap: 1027600384},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 500609097728, Avail: 20873912320},
			},
			Devices: []any{
				pci.Device{
					Bus:         "pci",
					Slot:        "0000:00:02.0",
					BusNumber:   0x0,
					DeviceClass: 0x300,
					VendorId:    0x8086,
					DeviceId:    0x9A49,
					SubvendorId: new(types.HexInt(0x1028)),
					SubdeviceId: new(types.HexInt(0xB20)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("TigerLake-LP GT2 [Iris Xe Graphics]"),
						SubvendorName: new("Dell"),
					},
				},
			},
		}, nil
	case "i7-2600k+arc-a580":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ht", "tm", "pbe", "syscall", "nx", "rdtscp", "lm", "constant_tsc", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "lahf_lm", "epb", "pti", "ssbd", "ibrs", "ibpb", "stibp", "tpr_shadow", "flexpriority", "ept", "vpid", "xsaveopt", "dtherm", "ida", "arat", "pln", "pts", "vnmi", "md_clear", "flush_l1d"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 11439935488, TotalSwap: 2147479552},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 58225459200, Avail: 16706797568},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x56A2,
					SubvendorId:          new(types.HexInt(0x1849)),
					SubdeviceId:          new(types.HexInt(0x6003)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("DG2 [Arc A580]"),
						SubvendorName: new("ASRock Incorporation"),
					},
					AdditionalProperties: map[string]string{
						"vram": "8096681984",
					},
				},
			},
		}, nil
	case "i7-10510U":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "sgx", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "mpx", "rdseed", "adx", "smap", "clflushopt", "intel_pt", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "vnmi", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 15908761600, TotalSwap: 4294963200},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1003736440832, Avail: 732326371328},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x9B41,
					SubvendorId:          new(types.HexInt(0x1028)),
					SubdeviceId:          new(types.HexInt(0x962)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("CometLake-U GT2 [UHD Graphics]"),
						SubvendorName: new("Dell"),
					},
				},
			},
		}, nil
	case "mustang":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "tme", "rdpid", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32876765184, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 502392942592, Avail: 431927586816},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:03:00.0",
					BusNumber:            0x3,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0xE20B,
					SubvendorId:          new(types.HexInt(0x8086)),
					SubdeviceId:          new(types.HexInt(0x1100)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Battlemage G21 [Arc B580]"),
						SubvendorName: new("Intel Corporation"),
					},
					AdditionalProperties: map[string]string{
						"vram": "12168933376",
					},
				},
			},
		}, nil
	case "raspberry-pi-5":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0B,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 8323276800, TotalSwap: 0},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 62270910464, Avail: 40666140672},
			},
		}, nil
	case "raspberry-pi-5+hailo-8":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:  "arm64",
					ImplementerId: 0x41,
					PartNumber:    0xD0B,
					Features:      []string{"fp", "asimd", "evtstrm", "aes", "pmull", "sha1", "sha2", "crc32", "atomics", "fphp", "asimdhp", "cpuid", "asimdrdm", "lrcpc", "dcpop", "asimddp"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 8317198336, TotalSwap: 1073737728},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 30777720832, Avail: 16953257984},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0001:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0xB40,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x1E60,
					DeviceId:             0x2864,
					SubvendorId:          new(types.HexInt(0x1E60)),
					SubdeviceId:          new(types.HexInt(0x2864)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Hailo Technologies Ltd."),
						DeviceName:    new("Hailo-8 AI Processor"),
						SubvendorName: new("Hailo Technologies Ltd."),
					},
				},
			},
		}, nil
	case "system76-addw4":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 32945307648, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 980105256960, Avail: 574739996672},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0xA788,
					SubvendorId:          new(types.HexInt(0x1558)),
					SubdeviceId:          new(types.HexInt(0x353)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Raptor Lake-S UHD Graphics"),
						SubvendorName: new("CLEVO/KAPOK Computer"),
					},
					AdditionalProperties: map[string]string{
						"vram": "31002615808",
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:01:00.0",
					BusNumber:            0x1,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x10DE,
					DeviceId:             0x28A1,
					SubvendorId:          new(types.HexInt(0x1558)),
					SubdeviceId:          new(types.HexInt(0x353)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("NVIDIA Corporation"),
						DeviceName:    new("AD107M [GeForce RTX 4050 Max-Q / Mobile]"),
						SubvendorName: new("CLEVO/KAPOK Computer"),
					},
					AdditionalProperties: map[string]string{
						"compute_capability": "8.9",
						"vram":               "6439305216",
					},
				},
			},
		}, nil
	case "xps13-7390":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "sgx", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "mpx", "rdseed", "adx", "smap", "clflushopt", "intel_pt", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "vnmi", "md_clear", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 15908761600, TotalSwap: 4294963200},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 1003736440832, Avail: 732343107584},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x9B41,
					SubvendorId:          new(types.HexInt(0x1028)),
					SubdeviceId:          new(types.HexInt(0x962)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("CometLake-U GT2 [UHD Graphics]"),
						SubvendorName: new("Dell"),
					},
					AdditionalProperties: map[string]string{
						"vram": "14482374656",
					},
				},
			},
		}, nil
	case "xps13-9350":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdt_a", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "lam", "wbnoinvd", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
				{
					Architecture:   "amd64",
					ManufacturerId: "GenuineIntel",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "dts", "acpi", "mmx", "fxsr", "sse", "sse2", "ss", "ht", "tm", "pbe", "syscall", "nx", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "art", "arch_perfmon", "pebs", "bts", "rep_good", "nopl", "xtopology", "nonstop_tsc", "cpuid", "aperfmperf", "tsc_known_freq", "pni", "pclmulqdq", "dtes64", "monitor", "ds_cpl", "vmx", "smx", "est", "tm2", "ssse3", "sdbg", "fma", "cx16", "xtpr", "pdcm", "pcid", "sse4_1", "sse4_2", "x2apic", "movbe", "popcnt", "tsc_deadline_timer", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "abm", "3dnowprefetch", "cpuid_fault", "epb", "intel_ppin", "ssbd", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "tpr_shadow", "flexpriority", "ept", "vpid", "ept_ad", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "rdt_a", "rdseed", "adx", "smap", "clflushopt", "clwb", "intel_pt", "sha_ni", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "split_lock_detect", "user_shstk", "avx_vnni", "lam", "wbnoinvd", "dtherm", "ida", "arat", "pln", "pts", "hwp", "hwp_notify", "hwp_act_window", "hwp_epp", "hwp_pkg_req", "hfi", "vnmi", "umip", "pku", "ospke", "waitpkg", "gfni", "vaes", "vpclmulqdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "fsrm", "md_clear", "serialize", "pconfig", "arch_lbr", "ibt", "flush_l1d", "arch_capabilities"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 33196507136, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 993185353728, Avail: 922389172224},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:02.0",
					BusNumber:            0x0,
					DeviceClass:          0x300,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x64A0,
					SubvendorId:          new(types.HexInt(0x1028)),
					SubdeviceId:          new(types.HexInt(0xCC9)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Lunar Lake [Intel Arc Graphics 130V / 140V]"),
						SubvendorName: new("Dell"),
					},
				},
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:00:0b.0",
					BusNumber:            0x0,
					DeviceClass:          0x1200,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x8086,
					DeviceId:             0x643E,
					SubvendorId:          new(types.HexInt(0x1028)),
					SubdeviceId:          new(types.HexInt(0xCC9)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Intel Corporation"),
						DeviceName:    new("Lunar Lake NPU"),
						SubvendorName: new("Dell"),
					},
				},
			},
		}, nil
	case "lenovo-thinkpad-p16s":
		return &machine.MachineInfo{
			Cpus: []cpu.CpuInfo{
				{
					Architecture:   "amd64",
					ManufacturerId: "AuthenticAMD",
					Flags:          []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clflush", "mmx", "fxsr", "sse", "sse2", "ht", "syscall", "nx", "mmxext", "fxsr_opt", "pdpe1gb", "rdtscp", "lm", "constant_tsc", "rep_good", "amd_lbr_v2", "nopl", "xtopology", "nonstop_tsc", "cpuid", "extd_apicid", "aperfmperf", "rapl", "pni", "pclmulqdq", "monitor", "ssse3", "fma", "cx16", "sse4_1", "sse4_2", "movbe", "popcnt", "aes", "xsave", "avx", "f16c", "rdrand", "lahf_lm", "cmp_legacy", "svm", "extapic", "cr8_legacy", "abm", "sse4a", "misalignsse", "3dnowprefetch", "osvw", "ibs", "skinit", "wdt", "tce", "topoext", "perfctr_core", "perfctr_nb", "bpext", "perfctr_llc", "mwaitx", "cpuid_fault", "cpb", "cat_l3", "cdp_l3", "hw_pstate", "ssbd", "mba", "perfmon_v2", "ibrs", "ibpb", "stibp", "ibrs_enhanced", "vmmcall", "fsgsbase", "tsc_adjust", "bmi1", "avx2", "smep", "bmi2", "erms", "invpcid", "cqm", "rdt_a", "avx512f", "avx512dq", "adx", "smap", "avx512ifma", "clflushopt", "clwb", "avx512cd", "sha_ni", "avx512bw", "avx512vl", "xsaveopt", "xsavec", "xgetbv1", "xsaves", "cqm_llc", "cqm_occup_llc", "cqm_mbm_total", "cqm_mbm_local", "user_shstk", "avx_vnni", "avx512_bf16", "clzero", "irperf", "xsaveerptr", "rdpru", "wbnoinvd", "cppc", "arat", "npt", "lbrv", "svm_lock", "nrip_save", "tsc_scale", "vmcb_clean", "flushbyasid", "decodeassists", "pausefilter", "pfthreshold", "avic", "v_vmsave_vmload", "vgif", "x2avic", "v_spec_ctrl", "vnmi", "avx512vbmi", "umip", "pku", "ospke", "avx512_vbmi2", "gfni", "vaes", "vpclmulqdq", "avx512_vnni", "avx512_bitalg", "avx512_vpopcntdq", "rdpid", "bus_lock_detect", "movdiri", "movdir64b", "overflow_recov", "succor", "smca", "fsrm", "avx512_vp2intersect", "flush_l1d", "amd_lbr_pmc_freeze"},
				},
			},
			Memory: memory.MemoryInfo{TotalRam: 92413652992, TotalSwap: 8589930496},
			Disk: map[string]disk.DirInfo{
				"/var/lib/snapd/snaps": {Total: 2011859927040, Avail: 521356337152},
			},
			Devices: []any{
				pci.Device{
					Bus:                  "pci",
					Slot:                 "0000:c4:00.0",
					BusNumber:            0xC4,
					DeviceClass:          0x380,
					ProgrammingInterface: new(uint8(0)),
					VendorId:             0x1002,
					DeviceId:             0x1114,
					SubvendorId:          new(types.HexInt(0x17AA)),
					SubdeviceId:          new(types.HexInt(0x512F)),
					FriendlyNames: pci.FriendlyNames{
						VendorName:    new("Advanced Micro Devices, Inc. [AMD/ATI]"),
						DeviceName:    new("Krackan [Radeon 840M / 860M Graphics]"),
						SubvendorName: new("Lenovo"),
					},
					AdditionalProperties: map[string]string{
						"microarchitecture": "gfx1152",
						"vram":              "8589934592",
					},
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("no machine fixture for %q", name)
	}
}
