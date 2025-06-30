package cpu

import (
	"log"
	"os"
	"testing"
)

var procCpuInfoTestFiles = map[string]string{
	"../../../test_data/devices/ampere-one-m-banshee-12/cpuinfo.txt":           arm64,
	"../../../test_data/devices/ampere-one-siryn/cpuinfo.txt":                  arm64,
	"../../../test_data/devices/ampere-one-x-banshee-8/cpuinfo.txt":            arm64,
	"../../../test_data/devices/hp-proliant-rl300-gen11-altra/cpuinfo.txt":     arm64,
	"../../../test_data/devices/hp-proliant-rl300-gen11-altra-max/cpuinfo.txt": arm64,
	"../../../test_data/devices/i7-2600k+arc-a580/cpuinfo.txt":                 amd64,
	"../../../test_data/devices/i7-10510U/cpuinfo.txt":                         amd64,
	"../../../test_data/devices/mustang/cpuinfo.txt":                           amd64,
	//"../../../test_data/devices/orange-pi-rv2/cpuinfo.txt":                     riscv64,
	"../../../test_data/devices/raspberry-pi-5/cpuinfo.txt":         arm64,
	"../../../test_data/devices/raspberry-pi-5+hailo-8/cpuinfo.txt": arm64,
	"../../../test_data/devices/xps13-7390/cpuinfo.txt":             amd64,
	"../../../test_data/devices/xps13-9350/cpuinfo.txt":             amd64,
}

func TestParseProcCpuInfo(t *testing.T) {

	for procCpuInfoFile, arch := range procCpuInfoTestFiles {
		t.Run(procCpuInfoFile, func(t *testing.T) {
			procCpuInfoBytes, err := os.ReadFile(procCpuInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			parsed, err := parseProcCpuInfo(string(procCpuInfoBytes), arch)
			if err != nil {
				t.Fatal(err)
			}

			for _, cpuInfo := range parsed {
				log.Printf("%+v", cpuInfo)
			}

		})
	}
}

func TestParseProcCpuInfoAmd64(t *testing.T) {
	cpuInfoData, err := os.ReadFile("../../../test_data/devices/xps13-7390/cpuinfo.txt")
	if err != nil {
		t.Fatal(err)
	}

	cpuInfos, err := parseProcCpuInfoAmd64(string(cpuInfoData))
	if err != nil {
		t.Fatal(err)
	}

	for _, cpuInfo := range cpuInfos {
		log.Printf("%+v", cpuInfo)
	}
}

func TestParseProcCpuInfoArm64(t *testing.T) {
	cpuInfoData, err := os.ReadFile("../../../test_data/devices/raspberry-pi-5/cpuinfo.txt")
	if err != nil {
		t.Fatal(err)
	}

	cpuInfos, err := parseProcCpuInfoArm64(string(cpuInfoData))
	if err != nil {
		t.Fatal(err)
	}

	for _, cpuInfo := range cpuInfos {
		log.Printf("%+v", cpuInfo)
	}
}
