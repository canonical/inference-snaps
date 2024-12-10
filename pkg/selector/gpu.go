package selector

import (
	"errors"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/gpu"
	"github.com/canonical/ml-snap-utils/pkg/types"
	"github.com/canonical/ml-snap-utils/pkg/utils"
)

func checkGpus(gpus []gpu.Gpu, stackDevice types.StackDevice) (float64, error) {
	for _, gpu := range gpus {
		result, err := gpuMatchesStack(gpu, stackDevice)
		if err != nil {
			return 0, err
		}
		if result {
			// At the first matching GPU we stop and return
			return 1, nil
		}
	}
	// If we get here, we checked all the GPUs and none were matches
	return 0, nil
}

// gpuMatchesStack checks if the GPU matches what is required by the stack definition.
// This is done as a filter, based on the fields in the stack definition.
// If the GPU from the hardware info passes all these filters, the GPU is a match.
func gpuMatchesStack(gpu gpu.Gpu, stackDevice types.StackDevice) (bool, error) {

	// If the stack has a Vendor ID requirement, check if the GPU's vendor matches
	// Vendor IDs are hex number strings, so do a case-insensitive compare
	if stackDevice.VendorId != nil && !strings.EqualFold(*stackDevice.VendorId, gpu.VendorId) {
		return false, nil
	}

	// If stack has a vram requirement, check if GPU has enough
	if stackDevice.VRam != nil {
		vramRequired, err := utils.StringToBytes(*stackDevice.VRam)
		if err != nil {
			return false, err
		}
		if vramAvailInterface, ok := gpu.Properties["vram"]; ok {
			vramAvailUInt, ok := vramAvailInterface.(uint64)
			if !ok {
				// hw info should list it as a uint64
				return false, errors.New("vram property is not uint64")
			}
			if vramAvailUInt < vramRequired {
				// Not enough vram
				return false, nil
			}
		} else {
			// Hardware Info does not list available vram
			return false, nil
		}
	}

	// TODO model id, compute capabilities

	// If we get here, all the filters have passed and the GPU is a match for this stack device
	return true, nil
}
