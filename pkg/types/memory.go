package types

type MemoryInfo struct {
	TotalRam  uint64 `json:"total_ram"`
	TotalSwap uint64 `json:"total_swap"`
}
