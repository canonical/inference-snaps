package pci

import (
	"errors"
	"fmt"
	"os"

	"github.com/canonical/stack-utils/pkg/constants"
	"github.com/canonical/stack-utils/pkg/hardware_info/pci/amd"
	"github.com/canonical/stack-utils/pkg/hardware_info/pci/intel"
	"github.com/canonical/stack-utils/pkg/hardware_info/pci/nvidia"
	"github.com/canonical/stack-utils/pkg/types"
	"github.com/jaypipes/pcidb"
)

var (
	pciDb                   *pcidb.PCIDB
	ErrorVendorNotSupported = errors.New("vendor not supported")
)

/*
Devices returns a slice of PciDevices that is detected on the current system and reported by lspci.
*/
func Devices(friendlyNames bool) ([]types.PciDevice, error) {

	hostLsPciData, err := hostLsPci()
	if err != nil {
		return nil, fmt.Errorf("error getting host lspci data: %v", err)
	}
	devices, err := ParseLsPci(hostLsPciData, friendlyNames)
	if err != nil {
		return nil, fmt.Errorf("error parsing lspci data: %v", err)
	}

	// Additional properties are obtained by running vendor specific tools on the host
	// Errors are not fatal, and are printed to stderr
	devices = addAdditionalProperties(devices)

	return devices, nil
}

func DevicesFromRawData(lspciData string, friendlyNames bool) ([]types.PciDevice, error) {
	devices, err := ParseLsPci(lspciData, friendlyNames)
	if err != nil {
		return nil, fmt.Errorf("error parsing lspci data: %v", err)
	}

	return devices, nil
}

/*
friendlyNames uses the numeric PCI ID fields to look up human-readable names for the device from the pci.id database.
*/
func friendlyNames(device types.PciDevice) (types.PciFriendlyNames, error) {
	var friendlyNames types.PciFriendlyNames

	if pciDb == nil {
		// Load pci.ids database if needed
		var err error
		pciDb, err = pcidb.New(pcidb.WithEnableNetworkFetch())
		if err != nil {
			return friendlyNames, fmt.Errorf("error opening pci database: %v", err)
		}
	}

	vendorIdString := fmt.Sprintf("%04x", device.VendorId)
	deviceIdString := fmt.Sprintf("%04x", device.DeviceId)

	subVendorIdString := ""
	if device.SubdeviceId != nil {
		subVendorIdString = fmt.Sprintf("%04x", *device.SubvendorId)
	}

	subDeviceIdString := ""
	if device.SubdeviceId != nil {
		subDeviceIdString = fmt.Sprintf("%04x", *device.SubdeviceId)
	}

	for _, vendor := range pciDb.Vendors {
		if vendor.ID == vendorIdString {
			vendorName := vendor.Name
			friendlyNames.VendorName = &vendorName

			for _, product := range vendor.Products {
				if product.ID == deviceIdString {
					productName := product.Name
					friendlyNames.DeviceName = &productName

					// Look up subDevice name from subsystem list
					if device.SubdeviceId != nil {
						for _, subSystem := range product.Subsystems {
							if subSystem.ID == subDeviceIdString {
								subSystemName := subSystem.Name
								friendlyNames.SubdeviceName = &subSystemName
							}
						}
					}
				}
			}
		}

		// Look up SubVendor name from main vendor list
		if device.SubvendorId != nil && vendor.ID == subVendorIdString {
			vendorName := vendor.Name
			friendlyNames.SubvendorName = &vendorName
		}
	}

	return friendlyNames, nil
}

/*
addAdditionalProperties returns devices with their AdditionalProperties field populated with device specific properties.
Additional properties are obtained by running vendor specific tools on the host system.
No error is returned as a failure to look up properties is considered non-fatal, and likely due to missing drivers.
Errors are instead logged to STDERR.
*/
func addAdditionalProperties(devices []types.PciDevice) []types.PciDevice {
	for i, device := range devices {
		properties, err := deviceAdditionalProperties(device)
		if err != nil {
			if errors.Is(err, ErrorVendorNotSupported) {
				// We do not log unsupported vendors, as that would be the majority of PCI devices
			} else {
				fmt.Fprintf(os.Stderr, "Error getting additional properties for pci device: %v\n", err)
			}
		}
		devices[i].AdditionalProperties = properties
	}

	return devices
}

func deviceAdditionalProperties(device types.PciDevice) (map[string]string, error) {
	var properties map[string]string
	var err error

	switch device.VendorId {
	case constants.PciVendorAmd:
		properties, err = amd.AdditionalProperties(device)
		if err != nil {
			return nil, fmt.Errorf("AMD: %v", err)
		}
	case constants.PciVendorNvidia:
		properties, err = nvidia.AdditionalProperties(device)
		if err != nil {
			return nil, fmt.Errorf("NVIDIA: %v", err)
		}
	case constants.PciVendorIntel:
		properties, err = intel.AdditionalProperties(device)
		if err != nil {
			return nil, fmt.Errorf("Intel: %v", err)
		}
	default:
		return nil, ErrorVendorNotSupported
	}

	return properties, nil
}
