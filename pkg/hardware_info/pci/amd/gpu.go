package amd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func gpuProperties(pciDevice types.PciDevice) (map[string]string, error) {
	return gpuPropertiesFromDir(pciDevice, "/")
}

func gpuPropertiesFromDir(pciDevice types.PciDevice, rootDir string) (map[string]string, error) {
	properties := make(map[string]string)

	vRamVal, err := vRam(pciDevice, rootDir)
	if err != nil {
		return nil, fmt.Errorf("error looking up vRAM: %v", err)
	}
	if vRamVal != nil {
		properties["vram"] = strconv.FormatUint(*vRamVal, 10)
	}
	gfxArchitecture, err := gfxArchitecture(pciDevice, rootDir)
	if err != nil {
		return nil, fmt.Errorf("error looking up gfx architecture: %v", err)
	}
	if len(gfxArchitecture) > 0 {
		properties["microarchitecture"] = gfxArchitecture
	}

	return properties, nil
}

func vRam(device types.PciDevice, rootDir string) (*uint64, error) {
	/*
		AMD vram is listed under /sys/bus/pci/devices/${pci_slot}/mem_info_vram_total

		ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_
		mem_info_gtt_total       mem_info_vis_vram_total  mem_info_vram_used
		mem_info_gtt_used        mem_info_vis_vram_used   mem_info_vram_vendor
		mem_info_preempt_used    mem_info_vram_total

		ubuntu@u-HP-EliteBook-845-G8-Notebook-PC:~$ cat /sys/bus/pci/devices/0000\:04\:00.0/mem_info_vram_total
		536870912
	*/
	data, err := os.ReadFile(filepath.Join(rootDir, "sys/bus/pci/devices", device.Slot, "mem_info_vram_total"))
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

func gfxArchitecture(device types.PciDevice, rootDir string) (string, error) {
	nodesDir := filepath.Join(rootDir, "sys/class/kfd/kfd/topology/nodes")
	files, err := os.ReadDir(nodesDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			propertiesPath := filepath.Join(nodesDir, file.Name(), "properties")
			data, err := os.ReadFile(propertiesPath)
			if err != nil {
				continue // skip this node if we can't read its properties
			}

			nodeMatchesDevice := false
			var nodeGfxTargetVersion string
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "drm_render_minor") {
					pciSlot, err := getAmdGpuPciSlot(line, rootDir)
					if err != nil {
						break
					}
					if pciSlot != device.Slot {
						break
					}
					nodeMatchesDevice = true
				} else if strings.HasPrefix(line, "gfx_target_version") {
					nodeGfxTargetVersion, err = parseGfxTargetVersion(line)
					if err != nil {
						break
					}
				}
			}

			if nodeMatchesDevice && len(nodeGfxTargetVersion) > 0 {
				return nodeGfxTargetVersion, nil
			}

		}
	}
	return "", fmt.Errorf("gfx_target_version not found for device with pci slot %s", device.Slot)
}

func getAmdGpuPciSlot(drmRenderMinor string, rootDir string) (string, error) {
	parts := strings.Split(drmRenderMinor, " ")
	if len(parts) == 2 {
		renderMinor := parts[1]
		pciSlotFull, err := filepath.EvalSymlinks(filepath.Join(rootDir, "sys/class/drm/renderD"+renderMinor, "device"))
		if err != nil {
			return "", err
		}
		pciSlot := strings.Split(string(pciSlotFull), "/")
		if len(pciSlot) == 0 {
			return "", fmt.Errorf("unexpected format for pci slot path: %s", pciSlotFull)
		}
		pciSlotStr := pciSlot[len(pciSlot)-1]
		return pciSlotStr, nil
	} else {
		return "", fmt.Errorf("unexpected format for drm_render_minor: %s", drmRenderMinor)
	}
}

func parseGfxTargetVersion(gfxTargetVersionLine string) (string, error) {
	parts := strings.Split(gfxTargetVersionLine, " ")
	if len(parts) == 2 {
		if parts[1] == "0" {
			return "", fmt.Errorf("gfx_target_version is invalid for this device")
		}
		gfxTargetVersion := parts[1]
		deviceLower := strings.ToLower(gfxTargetVersion)
		if len(deviceLower) < 6 {
			return "", fmt.Errorf("gfx_target_version has an unexpected format: %s", gfxTargetVersion)
		}

		majorInt, err := strconv.Atoi(deviceLower[0:2])
		if err != nil {
			return "", fmt.Errorf("error parsing major version from gfx_target_version: %v", err)
		}
		major := strconv.Itoa(majorInt)

		minorInt, err := strconv.Atoi(deviceLower[2:4])
		if err != nil {
			return "", fmt.Errorf("error parsing minor version from gfx_target_version: %v", err)
		}
		minor := strconv.Itoa(minorInt)

		revisionInt, err := strconv.Atoi(deviceLower[4:6])
		if err != nil {
			return "", fmt.Errorf("error parsing revision from gfx_target_version: %v", err)
		}
		revision := strconv.Itoa(revisionInt)

		arch := "gfx" + major + minor + revision
		return arch, nil
	}
	return "", fmt.Errorf("unexpected format for gfx_target_version: %s", gfxTargetVersionLine)
}
