package pci

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
)

func hostLsPci() (string, error) {
	// lspci -vmmnD
	out, err := exec.Command("lspci", "-vmmnD").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ParseLsPci(inputString string, includeFriendlyNames bool) ([]types.PciDevice, error) {
	var devices []types.PciDevice

	for _, section := range strings.Split(inputString, "\n\n") {
		// Ignore empty devices, e.g. extra blank line at end
		if section == "" {
			continue
		}
		var device types.PciDevice
		for _, line := range strings.Split(section, "\n") {
			key, value, _ := strings.Cut(line, ":\t")

			switch key {
			case "Slot":
				device.Slot = value
			case "Class":
				// e.g. 0x0300 for VGA controller
				if class, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.DeviceClass = types.HexInt(class)
				}
			case "Vendor":
				if vendor, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.VendorId = types.HexInt(vendor)
				}
			case "Device":
				if deviceId, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.DeviceId = types.HexInt(deviceId)
				}
			case "SVendor":
				if subVendorId, err := strconv.ParseUint(value, 16, 16); err == nil {
					subVendorIdUint16 := types.HexInt(subVendorId)
					device.SubvendorId = &subVendorIdUint16
				}
			case "SDevice":
				if subDeviceId, err := strconv.ParseUint(value, 16, 16); err == nil {
					subDeviceIdUint16 := types.HexInt(subDeviceId)
					device.SubdeviceId = &subDeviceIdUint16
				}
			case "ProgIf":
				// e.g. 0x02
				if progIf, err := strconv.ParseUint(value, 16, 8); err == nil {
					progIfUint8 := uint8(progIf)
					device.ProgrammingInterface = &progIfUint8
				}
			}

		}
		if includeFriendlyNames {
			friendlyNames, err := friendlyNames(device)
			if err != nil {
				// This is not a fatal error, so just logging it
				fmt.Fprintln(os.Stderr, "Error looking up friendly name:", err)
			} else {
				device.PciFriendlyNames = friendlyNames
			}
		}
		devices = append(devices, device)
	}

	return devices, nil
}
