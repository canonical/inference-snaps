package types

type DirStats struct {
	Total uint64 `json:"total" yaml:"total"`
	Avail uint64 `json:"avail" yaml:"avail"`
}
