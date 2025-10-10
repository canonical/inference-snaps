package engines

import (
	"fmt"
	"reflect"
	"slices"
)

func (device Device) validateBus(extraFields []string) error {
	switch device.Bus {
	case "pci":
		return device.validatePci(extraFields)
	case "usb":
		return device.validateUsb(extraFields)
	case "": // default to pci bus
		return device.validatePci(extraFields)
	default:
		return fmt.Errorf("invalid bus: %v", device.Bus)
	}
}

func (device Device) validateUsb(extraFields []string) error {
	return fmt.Errorf("usb device validation not implemented")
}

func (device Device) validatePci(extraFields []string) error {
	validFields := []string{
		"Type",
		"Bus",
		"VendorId",
		"DeviceId",
		"SnapConnections",
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
