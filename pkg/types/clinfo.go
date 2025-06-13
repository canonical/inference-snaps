package types

type Clinfo struct {
	Devices []struct {
		Online []struct {
			ClDeviceName              string `json:"CL_DEVICE_NAME"`
			ClDeviceVendor            string `json:"CL_DEVICE_VENDOR"`
			ClDeviceVendorID          int    `json:"CL_DEVICE_VENDOR_ID"`
			ClDeviceVersion           string `json:"CL_DEVICE_VERSION"`
			ClDevicePciBusInfoKhr     string `json:"CL_DEVICE_PCI_BUS_INFO_KHR"`
			ClDeviceGlobalMemSize     uint64 `json:"CL_DEVICE_GLOBAL_MEM_SIZE"`
			ClDeviceMaxMemAllocSize   int64  `json:"CL_DEVICE_MAX_MEM_ALLOC_SIZE"`
			ClDeviceHostUnifiedMemory bool   `json:"CL_DEVICE_HOST_UNIFIED_MEMORY"`
		} `json:"online"`
	} `json:"devices"`
}
