package engines

import (
	"fmt"
)

func (devices Devices) validate() error {
	for i, device := range devices.Allof {
		err := device.validate()
		if err != nil {
			return fmt.Errorf("invalid device: allof %d/%d: %v", i+1, len(devices.Allof), err)
		}
	}

	for i, device := range devices.Anyof {
		err := device.validate()
		if err != nil {
			return fmt.Errorf("invalid device: anyof %d/%d: %v", i+1, len(devices.Anyof), err)
		}
	}

	return nil
}

func (device Device) validate() error {
	switch device.Type {
	case "cpu":
		err := device.validateCpu()
		if err != nil {
			return fmt.Errorf("cpu: %v", err)
		}
	case "gpu":
		err := device.validateGpu()
		if err != nil {
			return fmt.Errorf("gpu: %v", err)
		}
	case "npu":
		err := device.validateNpu()
		if err != nil {
			return fmt.Errorf("npu: %v", err)
		}
	case "":
		err := device.validateTypelessDevice()
		if err != nil {
			return fmt.Errorf("typeless: %v", err)
		}
	default:
		return fmt.Errorf("invalid device type: %v", device.Type)
	}
	return nil
}

func (device Device) validateGpu() error {
	extraFields := []string{
		"VRam",
		"ComputeCapability",
	}

	err := device.validateBus(extraFields)
	if err != nil {
		return fmt.Errorf("gpu: %v", err)
	}
	return nil
}

func (device Device) validateNpu() error {
	err := device.validateBus(nil)
	if err != nil {
		return fmt.Errorf("npu: %v", err)
	}
	return nil
}

func (device Device) validateTypelessDevice() error {
	err := device.validateBus(nil)
	if err != nil {
		return fmt.Errorf("typeless device: %v", err)
	}
	return nil
}
