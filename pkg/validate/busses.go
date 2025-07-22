package validate

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/canonical/stack-utils/pkg/types"
)

func bus(device types.StackDevice, extraFields []string) error {
	switch device.Bus {
	case "pci":
		return pci(device, extraFields)
	case "usb":
		return usb(device, extraFields)
	case "": // default to pci bus
		return pci(device, extraFields)
	default:
		return fmt.Errorf("invalid bus: %v", device.Bus)
	}
}

func usb(device types.StackDevice, extraFields []string) error {
	return fmt.Errorf("usb device validation not implemented")
}

func pci(device types.StackDevice, extraFields []string) error {
	validFields := []string{
		"Type",
		"Bus",
		"VendorId",
		"DeviceId",
	}
	validFields = append(validFields, extraFields...)

	t := reflect.TypeOf(device)
	v := reflect.ValueOf(device)

	// Check fields with values against allow list
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			if !slices.Contains(validFields, fieldName) {
				return fmt.Errorf("pci device: invalid field: %s", fieldName)
			}
		}
	}

	return nil
}
