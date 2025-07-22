package validate

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/canonical/stack-utils/pkg/constants"
	"github.com/canonical/stack-utils/pkg/types"
)

func cpu(device types.StackDevice) error {
	if device.Architecture == nil {
		return fmt.Errorf("architecture field required")
	}

	switch *device.Architecture {
	case constants.Amd64:
		return cpuAmd64(device)
	case constants.Arm64:
		return cpuArm64(device)
	default:
		return fmt.Errorf("invalid architecture: %v", *device.Architecture)
	}
}

func cpuAmd64(device types.StackDevice) error {
	validFields := []string{
		"Type",
		"Architecture",
		"ManufacturerId",
		"Flags",
	}

	t := reflect.TypeOf(device)
	v := reflect.ValueOf(device)

	// Check fields with values against allow list
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			if !slices.Contains(validFields, fieldName) {
				return fmt.Errorf("cpu amd64: invalid field: %s", fieldName)
			}
		}
	}

	return nil
}

func cpuArm64(device types.StackDevice) error {
	validFields := []string{
		"Type",
		"Architecture",
		"ImplementerId",
		"PartNumber",
		"Features",
	}

	t := reflect.TypeOf(device)
	v := reflect.ValueOf(device)

	// Check fields with values against allow list
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.FieldByName(fieldName)
		if fieldValue.IsValid() && !fieldValue.IsZero() {
			if !slices.Contains(validFields, fieldName) {
				return fmt.Errorf("cpu arm64: invalid field: %s", fieldName)
			}
		}
	}

	return nil
}
