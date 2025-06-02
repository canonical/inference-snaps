package cpu

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/canonical/ml-snap-utils/pkg/types"
)

func hostLsCpu() ([]byte, error) {
	// lscpu --json --hierarchic
	out, err := exec.Command("lscpu", "--json", "--hierarchic").Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func parseLsCpu(input []byte) ([]types.CpuInfo, error) {
	var lsCpuJson lsCpuContainer
	err := json.Unmarshal(input, &lsCpuJson)
	if err != nil {
		return nil, err
	}

	var cpus []types.CpuInfo
	var architecture string
	var vendor string
	var modelName string

	for _, lsCpuObject := range lsCpuJson.LsCpu {
		label := lsCpuObject.Field
		value := lsCpuObject.Data

		switch label {
		case "Architecture:":
			architecture = value
		case "CPU(s):":
			// Not used as we calculate it ourselves per model
		case "Vendor ID:":
			vendor = value

			for _, vendorChild := range lsCpuObject.Children {
				switch vendorChild.Field {

				case "Model name:":
					modelName = vendorChild.Data

					var cpuInfo types.CpuInfo

					// Threads, cores, sockets and clusters are not always reported. Default to 1 of each.
					var threadsPerCore = 1
					var coresPerSocket = 1
					var coresPerCluster = 1
					var socketCount = 1
					var clusterCount = 1

					for _, modelNameChild := range vendorChild.Children {
						switch modelNameChild.Field {
						case "CPU family:":
							if familyId, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuInfo.FamilyId = &familyId
							}
						case "Model:":
							if modelId, err := strconv.Atoi(modelNameChild.Data); err == nil {
								cpuInfo.ModelId = modelId
							}
						case "Thread(s) per core:":
							if threads, err := strconv.Atoi(modelNameChild.Data); err == nil {
								threadsPerCore = threads
							}
						case "Core(s) per socket:":
							if cores, err := strconv.Atoi(modelNameChild.Data); err == nil {
								coresPerSocket = cores
							}
						case "Core(s) per cluster:":
							if cores, err := strconv.Atoi(modelNameChild.Data); err == nil {
								coresPerCluster = cores
							}
						case "Socket(s):":
							if sockets, err := strconv.Atoi(modelNameChild.Data); err == nil {
								socketCount = sockets
							}
						case "Cluster(s):":
							if clusters, err := strconv.Atoi(modelNameChild.Data); err == nil {
								clusterCount = clusters
							}
						case "CPU max MHz:":
							if maxFreq, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								cpuInfo.MaxFrequency = maxFreq
							}
						case "CPU min MHz:":
							if minFreq, err := strconv.ParseFloat(modelNameChild.Data, 64); err == nil {
								cpuInfo.MinFrequency = minFreq
							}
						case "BogoMIPS:":
							// Not used
						case "Flags:":
							flags := strings.Fields(modelNameChild.Data)
							cpuInfo.Flags = flags
						}
					}

					// Higher level data
					cpuInfo.Architecture = architecture
					cpuInfo.VendorId = vendor
					cpuInfo.ModelName = modelName

					// Calculate physical and logical cores
					cpuInfo.PhysicalCores = coresPerSocket * socketCount * coresPerCluster * clusterCount
					cpuInfo.LogicalCores = cpuInfo.PhysicalCores * threadsPerCore

					cpus = append(cpus, cpuInfo)
				}
			}
		}
	}

	return cpus, nil
}
