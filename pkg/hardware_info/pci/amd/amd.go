package amd

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
)

/*
AdditionalProperties returns device specific properties as a map[string]string.
No error is returned as a failure to look up properties is considered non-fatal, and likely due to missing drivers.
Any errors are logged to STDERR.
*/
func AdditionalProperties(pciDevice types.PciDevice) (map[string]string, error) {
	var properties map[string]string
	var err error

	// 00 01 - legacy VGA devices
	// 03 xx - display controllers
	if pciDevice.DeviceClass == 0x0001 || pciDevice.DeviceClass&0xFF00 == 0x0300 {
		properties, err = gpuProperties(pciDevice)
		if err != nil {
			return nil, fmt.Errorf("error getting gpu properties: %v", err)
		}
	}

	// Future: handle other AMD device classes

	return properties, nil
}
