package types

type PciDevice struct {
	Slot                 string  `json:"slot" yaml:"slot"`
	BusNumber            HexInt  `json:"bus_number" yaml:"bus-number"`
	DeviceClass          HexInt  `json:"device_class" yaml:"device-class"`
	ProgrammingInterface *uint8  `json:"programming_interface,omitempty" yaml:"programming-interface,omitempty"`
	VendorId             HexInt  `json:"vendor_id" yaml:"vendor-id"`
	DeviceId             HexInt  `json:"device_id" yaml:"device-id"`
	SubvendorId          *HexInt `json:"subvendor_id,omitempty" yaml:"subvendor-id,omitempty"`
	SubdeviceId          *HexInt `json:"subdevice_id,omitempty" yaml:"subdevice-id,omitempty"`
	PciFriendlyNames
	AdditionalProperties map[string]string `json:"additional_properties,omitempty" yaml:"additional-properties,omitempty"`
}

type PciFriendlyNames struct {
	VendorName    *string `json:"vendor_name,omitempty" yaml:"vendor-name,omitempty"`
	DeviceName    *string `json:"device_name,omitempty" yaml:"device-name,omitempty"`
	SubvendorName *string `json:"subvendor_name,omitempty" yaml:"subvendor-name,omitempty"`
	SubdeviceName *string `json:"subdevice_name,omitempty" yaml:"subdevice-name,omitempty"`
}
