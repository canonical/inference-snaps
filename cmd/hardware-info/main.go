package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/canonical/hardware-info/cpu"
	"github.com/canonical/hardware-info/disk"
	"github.com/canonical/hardware-info/memory"
)

func main() {
	var prettyOutput bool
	var fileOutput string

	flag.BoolVar(&prettyOutput, "pretty", false, "Output pretty json. Default is compact json.")
	flag.StringVar(&fileOutput, "file", "", "Output json to this file. Default output is to stdout.")
	flag.Parse()

	var hwInfo HwInfo

	memoryInfo, err := memory.Info()
	if err != nil {
		log.Println("Failed to get memory info:", err)
	}
	hwInfo.Memory = memoryInfo

	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Println("Failed to get CPU info:", err)
	}
	hwInfo.Cpu = cpuInfo

	diskInfo, err := disk.Info()
	if err != nil {
		log.Println("Failed to get disk info:", err)
	}
	hwInfo.Disk = diskInfo

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
