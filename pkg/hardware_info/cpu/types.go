package cpu

// ProcCpuInfo contains general information about a system CPU found in /proc/cpuinfo.
type ProcCpuInfo struct {
	Processor    int64 // %d - kernel defines it as long long
	Architecture string

	// amd64
	ManufacturerId string
	BrandString    string
	Flags          []string

	// arm64
	ModelName     *string  // %s
	BogoMips      float64  // %lu.%02lu
	Features      []string // space separated strings
	ImplementerId uint64   // 0x%02x
	//Architecture  uint64   // constant int
	Variant    uint64 // 0x%x
	PartNumber uint64 // 0x%03x
	Revision   uint64 // %d
}
