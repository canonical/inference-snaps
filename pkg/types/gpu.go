package types

type Gpu struct {
	Bus               string  `json:"bus"`
	VendorId          string  `json:"vendor_id"`
	VendorName        *string `json:"vendor_name,omitempty"`
	DeviceId          string  `json:"device_id"`
	DeviceName        *string `json:"device_name,omitempty"`
	SubvendorId       *string `json:"subvendor_id,omitempty"`
	SubvendorName     *string `json:"subvendor_name,omitempty"`
	SubdeviceId       *string `json:"subdevice_id,omitempty"`
	SubdeviceName     *string `json:"subdevice_name,omitempty"`
	VRam              *uint64 `json:"vram,omitempty"`
	ComputeCapability *string `json:"compute_capability,omitempty"`
}
