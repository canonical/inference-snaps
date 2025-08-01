package amd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
)

func gpuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice)
	if err != nil {
		return nil, fmt.Errorf("error looking up vRAM: %v", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}

	return properties, nil
}

func vRam(device types.PciDevice) (*uint64, error) {
	/*
		AMD vram is listed under /sys/bus/pci/devices/${pci_slot}/mem_info_vram_total

		ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_
		mem_info_gtt_total       mem_info_vis_vram_total  mem_info_vram_used
		mem_info_gtt_used        mem_info_vis_vram_used   mem_info_vram_vendor
		mem_info_preempt_used    mem_info_vram_total

		ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_vram_total
		536870912
	*/
	data, err := os.ReadFile("/sys/bus/pci/devices/" + device.Slot + "/mem_info_vram_total")
	if err != nil {
		return nil, err
	}
	dataStr := string(data)
	dataStr = strings.TrimSpace(dataStr) // value in file ends in \n
	vram, err := strconv.ParseUint(dataStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return &vram, nil
}
