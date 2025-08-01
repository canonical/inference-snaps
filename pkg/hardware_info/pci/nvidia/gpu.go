package nvidia

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
)

func gpuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice)
	if err != nil {
		return nil, fmt.Errorf("error looking up vRAM: %v", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}

	ccVal, err := computeCapability(pciDevice)
	if err != nil {
		return nil, fmt.Errorf("error looking up compute capability: %v", err)
	}
	if ccVal != nil {
		properties["compute_capability"] = *ccVal
	}

	return properties, nil
}

func vRam(device types.PciDevice) (*uint64, error) {
	/*
		Nvidia: LANG=C nvidia-smi --query-gpu=memory.total --format=csv,noheader,nounits

		$ nvidia-smi --id=00000000:01:00.0 --query-gpu=memory.total --format=csv,noheader
		4096 MiB
		$ nvidia-smi --id=00000000:02:00.0 --query-gpu=memory.total --format=csv,noheader
		No devices were found
	*/
	output, err := nvidiaSmi("--id="+device.Slot, "--query-gpu=memory.total", "--format=csv,noheader")
	if err != nil {
		return nil, fmt.Errorf("error executing nvidia-smi: %s", err)
	}

	valueStr, unit, hasUnit := strings.Cut(*output, " ")
	vramValue, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return nil, err
	}

	if hasUnit {
		switch unit {
		case "KiB":
			vramValue = vramValue * 1024
		case "MiB":
			vramValue = vramValue * 1024 * 1024
		case "GiB":
			vramValue = vramValue * 1024 * 1024 * 1024
		}
	}

	return &vramValue, nil

}

func computeCapability(device types.PciDevice) (*string, error) {
	// nvidia-smi --query-gpu=compute_cap --format=csv,noheader
	output, err := nvidiaSmi("--id="+device.Slot, "--query-gpu=compute_cap", "--format=csv,noheader")
	if err != nil {
		return nil, fmt.Errorf("error executing nvidia-smi: %s", err)
	}

	return output, nil
}

func nvidiaSmi(args ...string) (*string, error) {
	cmd := exec.Command("nvidia-smi", args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LANG=C")
	output, err := cmd.Output()
	if err != nil {
		if len(output) == 0 {
			return nil, err
		} else {
			// nvidia-smi writes error messages to stdout
			return nil, fmt.Errorf("%s: %s", err, bytes.TrimSpace(output))
		}
	}

	strOutput := string(bytes.TrimSpace(output))
	return &strOutput, nil
}
