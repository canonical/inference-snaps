package memory

type MemoryInfo struct {
	RamTotal  uint64 `json:"ram_total"`
	SwapTotal uint64 `json:"swap_total"`
}
