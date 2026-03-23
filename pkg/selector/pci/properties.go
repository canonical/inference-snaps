package pci

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

func checkProperties(device engines.Device, pciDevice types.PciDevice) (int, error) {
	extraScore := 0

	// vram
	if device.VRam != nil {
		err := checkVram(device, pciDevice)
		if err != nil {
			return 0, err
		}
		extraScore += weights.GpuVRam
	}

	// microarchitecture
	if device.Microarchitecture != nil {
		err := checkMicroarchitecture(*device.Microarchitecture, pciDevice)
		if err != nil {
			return 0, err
		}
		extraScore += weights.GpuMicroarchitecture
	}
	// TODO compute-capability

	return extraScore, nil
}

func checkVram(device engines.Device, pciDevice types.PciDevice) error {
	vramRequired, err := utils.StringToBytes(*device.VRam)
	if err != nil {
		return err
	}
	if vram, ok := pciDevice.AdditionalProperties["vram"]; ok {
		vramAvailable, err := utils.StringToBytes(vram)
		if err != nil {
			return fmt.Errorf("error parsing vRAM: %v", err)
		}
		if vramAvailable >= vramRequired {
			return nil
		} else {
			return fmt.Errorf("not enough vRAM: %d", vramAvailable)
		}
	} else {
		// Hardware Info does not list available vram
		return fmt.Errorf("unable to detect vRAM")
	}
}

func checkMicroarchitecture(microArchRequired string, pciDevice types.PciDevice) error {
	if microArch, ok := pciDevice.AdditionalProperties["microarchitecture"]; ok {
		if microArch == microArchRequired {
			return nil
		} else {
			return fmt.Errorf("microarchitecture does not match: %s", microArch)
		}
	} else {
		// Hardware Info does not list available microarchitecture
		return fmt.Errorf("unable to detect microarchitecture")
	}
}
