package cpu

func Info() (*CpuInfo, error) {
	hostLsCpu, err := hostLsCpu()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := parseLsCpu(hostLsCpu)
	if err != nil {
		return nil, err
	}

	return cpuInfo, err
}
