package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/cpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/disk"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/gpu"
	"github.com/canonical/ml-snap-utils/pkg/hardware_info/memory"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func main() {
	var prettyOutput bool
	var friendlyNames bool
	var fileOutput string

	flag.BoolVar(&prettyOutput, "pretty", false, "Output pretty json. Default is compact json.")
	flag.BoolVar(&friendlyNames, "friendly", false, "Include human readable names for devices.")
	flag.StringVar(&fileOutput, "file", "", "Output json to this file. Default output is to stdout.")
	flag.Parse()

	var hwInfo types.HwInfo

	memoryInfo, err := memory.Info()
	if err != nil {
		log.Fatalf("Failed to get memory info: %s", err)
	}
	hwInfo.Memory = memoryInfo

	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Fatalf("Failed to get CPU info: %s", err)
	}
	hwInfo.Cpu = cpuInfo

	diskInfo, err := disk.Info()
	if err != nil {
		log.Fatalf("Failed to get disk info: %s", err)
	}
	hwInfo.Disk = diskInfo

	gpuInfo, err := gpu.Info(friendlyNames)
	if err != nil {
		log.Fatalf("Failed to get GPU info: %s", err)
	}
	hwInfo.Gpus = gpuInfo

	var jsonString []byte
	if prettyOutput {
		jsonString, err = json.MarshalIndent(hwInfo, "", "  ")
		if err != nil {
			log.Fatal("Failed to marshal to JSON:", err)
		}
	} else {
		jsonString, err = json.Marshal(hwInfo)
		if err != nil {
			log.Fatal("Failed to marshal to JSON:", err)
		}
	}

	if fileOutput != "" {
		err = os.WriteFile(fileOutput, jsonString, 0644)
		if err != nil {
			log.Fatal("Failed to write to file:", err)
		}
	} else {
		fmt.Println(string(jsonString))
	}
}
