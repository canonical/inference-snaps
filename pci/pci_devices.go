package pci

import (
	"fmt"
	"github.com/jaypipes/pcidb"
)

var (
	pciDb *pcidb.PCIDB
)

func PciDevices(friendlyNames bool) ([]Device, error) {

	hostLsPci, err := hostLsPci()
	if err != nil {
		return nil, err
	}
	devices, err := ParseLsPci(hostLsPci, friendlyNames)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func lookupFriendlyNames(device Device) (FriendlyNames, error) {
	var friendlyNames FriendlyNames

	if pciDb == nil {
		// Load pci.ids database if needed
		var err error
		pciDb, err = pcidb.New()
		if err != nil {
			return friendlyNames, err
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
