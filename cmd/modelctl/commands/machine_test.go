package commands

import (
	"os"
	"testing"

	"github.com/canonical/lscompute/pkg/machine"
)

func Example_machineCommand_printMachineJson() {
	cmd := machineCommand{format: "json"}
	info, err := machineFixture("dummy-machine")
	if err != nil {
		panic(err)
	}
	machineDetails := NewMachineDetails(info)
	if err := cmd.printMachineJson(*machineDetails); err != nil {
		panic(err)
	}

	// Output:
	// {
	//   "cpus": [
	//     {
	//       "architecture": "amd64",
	//       "manufacturer-id": "GenuineIntel",
	//       "flags": [
	//         "fpu",
	//         "vme",
	//         "de"
	//       ]
	//     }
	//   ],
	//   "memory": {
	//     "total-ram": 67012501504,
	//     "total-swap": 0
	//   },
	//   "disks": [
	//     {
	//       "path": "/var/lib/snapd/snaps",
	//       "total": 1006451294208,
	//       "avail": 943543738368
	//     }
	//   ],
	//   "devices": [
	//     {
	//       "bus": "pci",
	//       "slot": "0000:00:00.0",
	//       "bus-number": "0x0",
	//       "device-class": "0x600",
	//       "vendor-id": "0x8086",
	//       "device-id": "0x4637",
	//       "subvendor-id": "0x103C",
	//       "subdevice-id": "0x89C6",
	//       "vendor-name": "Intel Corporation",
	//       "subvendor-name": "Hewlett-Packard Company"
	//     }
	//   ]
	// }

}

func Example_machineCommand_printMachineYaml() {
	cmd := machineCommand{format: "yaml"}
	info, err := machineFixture("dummy-machine")
	if err != nil {
		panic(err)
	}
	machineDetails := NewMachineDetails(info)
	if err := cmd.printMachineYaml(*machineDetails); err != nil {
		panic(err)
	}

	// Output:
	// cpus:
	//     - architecture: amd64
	//       manufacturer-id: GenuineIntel
	//       flags: [fpu, vme, de]
	// memory:
	//     total-ram: 62.4G
	//     total-swap: "0"
	// disks:
	//     - path: /var/lib/snapd/snaps
	//       total: 937.3G
	//       avail: 878.7G
	// devices:
	//     - bus: pci
	//       slot: '0000:00:00.0'
	//       bus-number: "0x0"
	//       device-class: "0x600"
	//       vendor-id: "0x8086"
	//       device-id: "0x4637"
	//       subvendor-id: "0x103C"
	//       subdevice-id: "0x89C6"
	//       vendor-name: Intel Corporation
	//       subvendor-name: Hewlett-Packard Company
}

func Test_newPciAdditionalDeviceProperties_invalidVRAM(t *testing.T) {
	stderr, err := os.CreateTemp(t.TempDir(), "stderr-*")
	if err != nil {
		t.Fatal(err)
	}
	originalStderr := os.Stderr
	os.Stderr = stderr
	t.Cleanup(func() {
		os.Stderr = originalStderr
		stderr.Close()
	})

	properties := newPciAdditionalDeviceProperties(map[string]any{"vram": "8G"})
	if properties != nil {
		t.Errorf("expected nil properties for invalid VRAM, got %+v", properties)
	}

	output, err := os.ReadFile(stderr.Name())
	if err != nil {
		t.Fatal(err)
	}
	want := "Warning: it was not possible to parse VRAM value: expected uint64, got string (8G).\n"
	if string(output) != want {
		t.Errorf("warning = %q, want %q", output, want)
	}
}

func Test_printMachine_unknownFormat(t *testing.T) {
	cmd := machineCommand{format: "xml"}
	info := &machine.Machine{}

	err := cmd.printMachine(*NewMachineDetails(info))
	if err == nil || err.Error() != `unknown format "xml"` {
		t.Errorf("expected error 'unknown format \"xml\"', got %v", err)
	}
}
