package types

type MemoryInfo struct {
	TotalRam  uint64 `json:"total_ram" yaml:"total-ram"`
	TotalSwap uint64 `json:"total_swap" yaml:"total-swap"`
}
