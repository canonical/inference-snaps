package cpu

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

func hostLsCpu() ([]byte, error) {
	out, err := exec.Command("lscpu", "--json", "--hierarchic").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func parseLsCpu(input []byte) (*CpuInfo, error) {
	cpuInfo := CpuInfo{}

	var lsCpuJson LsCpuContainer
	err := json.Unmarshal(input, &lsCpuJson)
	if err != nil {
		return nil, err
	}

	for _, lsCpuObject := range lsCpuJson.LsCpu {
		label := lsCpuObject.Field
		value := lsCpuObject.Data

		switch label {
		case "Architecture:":
			cpuInfo.Architecture = value
		case "CPU(s):":
			if cpuCount, err := strconv.Atoi(value); err == nil {
				cpuInfo.CpuCount = cpuCount
			}
		case "Vendor ID:":
			cpuInfo.Vendor = value

			for _, vendorChild := range lsCpuObject.Children {
				switch vendorChild.Field {

				case "Model name:":
					cpuModel := Model{Name: value}
					cpuModel.Name = vendorChild.Data

					for _, modelNameChild := range vendorChild.Children {
						switch modelNameChild.Field {
						case "CPU family:":
							if familyId, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.Family = &familyId
							}
						case "Model:":
							if modelId, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.Id = modelId
							}
						case "Thread(s) per core:":
							if threads, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.ThreadsPerCore = &threads
							}
						case "Core(s) per socket:":
							if cores, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.CoresPerSocket = &cores
							}
						case "Core(s) per cluster:":
							if cores, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.CoresPerCluster = &cores
							}
						case "Socket(s):":
							if sockets, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.Sockets = &sockets
							}
						case "Cluster(s):":
							if clusters, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuModel.Clusters = &clusters
							}
						case "CPU max MHz:":
							if maxFreq, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								cpuModel.MaxFreq = maxFreq
							}
						case "CPU min MHz:":
							if minFreq, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								cpuModel.MinFreq = minFreq
							}
						case "BogoMIPS:":
							if bogoMips, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								cpuModel.BogoMips = bogoMips
							}
						case "Flags:":
							flags := strings.Fields(modelNameChild.Data)
							cpuModel.Flags = flags
						}
					}
					cpuInfo.Models = append(cpuInfo.Models, cpuModel)
				}
			}
		}
	}

	return &cpuInfo, nil
}
