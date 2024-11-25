package common

type HwInfo struct {
	Cpu    *CpuInfo             `json:"cpu,omitempty"`
	Memory *MemoryInfo          `json:"memory,omitempty"`
	Disk   map[string]*DirStats `json:"disk,omitempty"`
	Gpus   []Gpu                `json:"gpu,omitempty"`
}

type CpuInfo struct {
	Architecture string     `json:"architecture"`
	CpuCount     int        `json:"cpu_count"`
	Vendor       string     `json:"vendor"`
	Models       []CpuModel `json:"models"`
}

type CpuModel struct {
	Name            string `json:"name"`
	Family          *int   `json:"family,omitempty"`
	Id              int    `json:"id"`
	ThreadsPerCore  *int   `json:"threads_per_core,omitempty"`
	Sockets         *int   `json:"sockets,omitempty"`
	CoresPerSocket  *int   `json:"cores_per_socket,omitempty"`
	Clusters        *int   `json:"clusters,omitempty"`
	CoresPerCluster *int   `json:"cores_per_cluster,omitempty"`
	//CpuCount int // = sockets * cores-per-socket * clusters * cores-per-cluster * threads-per-core
	MaxFreq  float64  `json:"max_freq"`
	MinFreq  float64  `json:"min_freq"`
	BogoMips float64  `json:"bogo_mips"`
	Flags    []string `json:"flags"`
}

type MemoryInfo struct {
	RamTotal  uint64 `json:"ram_total"`
	SwapTotal uint64 `json:"swap_total"`
}

type DirStats struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Avail uint64 `json:"avail"`
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
