package types

type CpuInfo struct {
	Architecture string `json:"architecture" yaml:"architecture"`

	// amd64
	ManufacturerId string   `json:"manufacturer-id,omitempty" yaml:"manufacturer-id,omitempty"`
	Flags          []string `json:"flags,omitempty" yaml:"flags,omitempty"`

	// arm64
	ImplementerId HexInt   `json:"implementer-id,omitempty" yaml:"implementer-id,omitempty"`
	PartNumber    HexInt   `json:"part-number,omitempty" yaml:"part-number,omitempty"`
	Features      []string `json:"features,omitempty" yaml:"features,omitempty"`
}
