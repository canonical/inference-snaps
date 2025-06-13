package types

type PciDevice struct {
	Slot                 string  `json:"slot"`
	DeviceClass          HexInt  `json:"device_class"`
	ProgrammingInterface *uint8  `json:"programming_interface,omitempty"`
	VendorId             HexInt  `json:"vendor_id"`
	DeviceId             HexInt  `json:"device_id"`
	SubvendorId          *HexInt `json:"subvendor_id,omitempty"`
	SubdeviceId          *HexInt `json:"subdevice_id,omitempty"`
	PciFriendlyNames
	AdditionalProperties map[string]string `json:"additional_properties,omitempty"`
}

type PciFriendlyNames struct {
	VendorName    *string `json:"vendor_name,omitempty"`
	DeviceName    *string `json:"device_name,omitempty"`
	SubvendorName *string `json:"subvendor_name,omitempty"`
	SubdeviceName *string `json:"subdevice_name,omitempty"`
}
