package pci

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/canonical/lscompute/pkg/machine/device/pci"
)

func checkProperties(manifestDevice engines.Device, hostPciDevice pci.Device) (int, error) {
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
	// compute-capability
	if manifestDevice.ComputeCapability != nil {
		err := checkComputeCapability(*manifestDevice.ComputeCapability, hostPciDevice)
		if err != nil {
			return 0, fmt.Errorf("checking compute-capability: %v", err)
		}
		extraScore += weights.GpuComputeCapability
	}

	return extraScore, nil
}

func checkVram(manifestDevice engines.Device, hostPciDevice pci.Device) error {
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

func checkMicroarchitecture(microArchRequired string, hostPciDevice pci.Device) error {
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

// checkComputeCapability compares the reported compute capability with a semver version range (e.g. ">=6.0, <7.0").
func checkComputeCapability(constraintRequired string, hostPciDevice pci.Device) error {
	constraint, err := semver.NewConstraint(constraintRequired)
	if err != nil {
		return fmt.Errorf("parsing constraint %q: %v", constraintRequired, err)
	}

	if cc, ok := hostPciDevice.AdditionalProperties["compute-capability"]; ok {
		version, err := semver.NewVersion(cc)
		if err != nil {
			return fmt.Errorf("parsing compute-capability %q: %v", cc, err)
		}
		if constraint.Check(version) {
			return nil
		}
		return fmt.Errorf("compute-capability %s does not satisfy %s", cc, constraintRequired)
	}

	// Hardware Info does not list a compute capability
	return fmt.Errorf("compute-capability not reported")
}
