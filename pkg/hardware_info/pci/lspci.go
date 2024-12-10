package pci

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func hostLsPci() ([]byte, error) {
	out, err := exec.Command("lspci", "-vmmnD").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ParseLsPci(input []byte, includeFriendlyNames bool) ([]Device, error) {
	var devices []Device

	inputString := string(input)
	for _, section := range strings.Split(inputString, "\n\n") {
		var device Device
		for _, line := range strings.Split(section, "\n") {
			key, value, _ := strings.Cut(line, ":\t")

			switch key {
			case "Slot":
				device.Slot = value
			case "Class":
				// e.g. 0x0300 for VGA controller
				if class, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.DeviceClass = uint16(class)
				}
			case "Vendor":
				if vendor, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.VendorId = uint16(vendor)
				}
			case "Device":
				if deviceId, err := strconv.ParseUint(value, 16, 16); err == nil {
					device.DeviceId = uint16(deviceId)
				}
			case "SVendor":
				if subVendorId, err := strconv.ParseUint(value, 16, 16); err == nil {
					subVendorIdUint16 := uint16(subVendorId)
					device.SubvendorId = &subVendorIdUint16
				}
			case "SDevice":
				if subDeviceId, err := strconv.ParseUint(value, 16, 16); err == nil {
					subDeviceIdUint16 := uint16(subDeviceId)
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
			friendlyNames, err := lookupFriendlyNames(device)
			if err != nil {
				// This is not a fatal error, so just logging it
				log.Printf("Error looking up friendly name: %v", err)
			} else {
				device.FriendlyNames = friendlyNames
			}
		}
		devices = append(devices, device)
	}

	return devices, nil
}
