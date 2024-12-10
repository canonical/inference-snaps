package gpu

type Display struct {
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
}

type Gpu struct {
	VendorId      string                 `json:"vendor_id"`
	VendorName    *string                `json:"vendor_name,omitempty"`
	DeviceId      string                 `json:"device_id"`
	DeviceName    *string                `json:"device_name,omitempty"`
	SubvendorId   *string                `json:"subvendor_id,omitempty"`
	SubvendorName *string                `json:"subvendor_name,omitempty"`
	SubdeviceId   *string                `json:"subdevice_id,omitempty"`
	SubdeviceName *string                `json:"subdevice_name,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}
