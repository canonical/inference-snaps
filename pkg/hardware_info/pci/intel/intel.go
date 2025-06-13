package intel

import (
	"github.com/canonical/ml-snap-utils/pkg/types"
)

/*
AdditionalProperties returns device specific properties as a map[string]string.
No error is returned as a failure to look up properties is considered non-fatal, and likely due to missing drivers.
Any errors are logged to STDERR.
*/
func AdditionalProperties(pciDevice types.PciDevice) map[string]string {
	var properties map[string]string

	// 00 01 - legacy VGA devices
	// 03 xx - display controllers
	if pciDevice.DeviceClass == 0x0001 || pciDevice.DeviceClass&0xFF00 == 0x0300 {
		properties = gpuProperties(pciDevice)
	}

	// Future: handle other Intel device classes like NPUs

	return properties
}
