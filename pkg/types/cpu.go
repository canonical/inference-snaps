package types

type CpuInfo struct {
	Architecture string `json:"architecture"`

	// amd64
	ManufacturerId string   `json:"manufacturer_id,omitempty"`
	Flags          []string `json:"flags,omitempty"`

	// arm64
	ImplementerId HexInt `json:"implementer_id,omitempty"`
	PartNumber    HexInt `json:"part_number,omitempty"`
}
