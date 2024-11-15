package cpu

type LsCpuContainer struct {
	LsCpu []LsCpuObject `json:"lscpu"`
}

type LsCpuObject struct {
	Field    string        `json:"field"`
	Data     string        `json:"data"`
	Children []LsCpuObject `json:"children"`
}

type CpuInfo struct {
	Architecture string  `json:"architecture"`
	CpuCount     int     `json:"cpu_count"`
	Vendor       string  `json:"vendor"`
	Models       []Model `json:"models"`
}

type Model struct {
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
