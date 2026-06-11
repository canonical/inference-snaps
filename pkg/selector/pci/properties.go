package pci

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/v2/pkg/types"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
)

func checkProperties(manifestDevice engines.Device, hostPciDevice types.PciDevice) (int, error) {
	extraScore := 0

	// vram
	if manifestDevice.VRam != nil {
		err := checkVram(manifestDevice, hostPciDevice)
		if err != nil {
			return 0, fmt.Errorf("checking vram: %v", err)
		}
		extraScore += weights.GpuVRam
	}

	// microarchitecture
	if manifestDevice.Microarchitecture != nil {
		err := checkMicroarchitecture(*manifestDevice.Microarchitecture, hostPciDevice)
		if err != nil {
			return 0, fmt.Errorf("checking microarchitecture: %v", err)
		}
		extraScore += weights.GpuMicroarchitecture
	}
	// TODO compute-capability

	return extraScore, nil
}

func checkVram(manifestDevice engines.Device, hostPciDevice types.PciDevice) error {
	vramRequired, err := utils.StringToBytes(*manifestDevice.VRam)
	if err != nil {
		return err
	}
	if vram, ok := hostPciDevice.AdditionalProperties["vram"]; ok {
		vramAvailable, err := utils.StringToBytes(vram)
		if err != nil {
			return fmt.Errorf("parsing vram: %v", err)
		}
		if vramAvailable >= vramRequired {
			return nil
		} else {
			return fmt.Errorf("not enough vram: %d", vramAvailable)
		}
	} else {
		// Hardware Info does not list available vram
		return fmt.Errorf("vram not reported")
	}
}

func checkMicroarchitecture(microArchRequired string, hostPciDevice types.PciDevice) error {
	if microArch, ok := hostPciDevice.AdditionalProperties["microarchitecture"]; ok {
		if microArch == microArchRequired {
			return nil
		} else {
			return fmt.Errorf("microarchitecture does not match: %s", microArch)
		}
	} else {
		// Hardware Info does not list available microarchitecture
		return fmt.Errorf("microarchitecture not reported")
	}
}
