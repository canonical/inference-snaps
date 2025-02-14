package cpu

type lsCpuContainer struct {
	LsCpu []lsCpuObject `json:"lscpu"`
}

type lsCpuObject struct {
	Field    string        `json:"field"`
	Data     string        `json:"data"`
	Children []lsCpuObject `json:"children"`
}
