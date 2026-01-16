package types

type PciDevice struct {
	Score                int     `json:"score,omitempty" yaml:"score,omitempty"`
	Slot                 string  `json:"slot" yaml:"slot"`
	BusNumber            HexInt  `json:"bus-number" yaml:"bus-number"`
	DeviceClass          HexInt  `json:"device-class" yaml:"device-class"`
	ProgrammingInterface *uint8  `json:"programming-interface,omitempty" yaml:"programming-interface,omitempty"`
	VendorId             HexInt  `json:"vendor-id" yaml:"vendor-id"`
	DeviceId             HexInt  `json:"device-id" yaml:"device-id"`
	SubvendorId          *HexInt `json:"subvendor-id,omitempty" yaml:"subvendor-id,omitempty"`
	SubdeviceId          *HexInt `json:"subdevice-id,omitempty" yaml:"subdevice-id,omitempty"`
	PciFriendlyNames     `yaml:",inline"`
	AdditionalProperties map[string]string `json:"additional-properties,omitempty" yaml:"additional-properties,omitempty"`
}

type PciFriendlyNames struct {
	VendorName    *string `json:"vendor-name,omitempty" yaml:"vendor-name,omitempty"`
	DeviceName    *string `json:"device-name,omitempty" yaml:"device-name,omitempty"`
	SubvendorName *string `json:"subvendor-name,omitempty" yaml:"subvendor-name,omitempty"`
	SubdeviceName *string `json:"subdevice-name,omitempty" yaml:"subdevice-name,omitempty"`
}
