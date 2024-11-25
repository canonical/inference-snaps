package pci

type Device struct {
	Slot                 string  `json:"slot"`
	DeviceClass          uint16  `json:"device_class"`
	ProgrammingInterface *uint8  `json:"programming_interface"`
	VendorId             uint16  `json:"vendor_id"`
	DeviceId             uint16  `json:"device_id"`
	SubvendorId          *uint16 `json:"subvendor_id,omitempty"`
	SubdeviceId          *uint16 `json:"subdevice_id,omitempty"`
	FriendlyNames
}

type FriendlyNames struct {
	VendorName    *string `json:"vendor_name,omitempty"`
	DeviceName    *string `json:"device_name,omitempty"`
	SubvendorName *string `json:"subvendor_name,omitempty"`
	SubdeviceName *string `json:"subdevice_name,omitempty"`
}
