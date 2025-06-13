package cpu

import (
	"log"
	"os"
	"testing"
)

var procCpuInfoTestFiles = map[string]string{
	"../../../test_data/cpuinfo/ampere-altra.txt":            arm64,
	"../../../test_data/cpuinfo/ampere-one-m-banshee-12.txt": arm64,
	"../../../test_data/cpuinfo/ampere-one-siryn.txt":        arm64,
	"../../../test_data/cpuinfo/ampere-one-x-banshee-8.txt":  arm64,
	"../../../test_data/cpuinfo/raspberry-pi-5.txt":          arm64,
	"../../../test_data/cpuinfo/xps13-7390.txt":              amd64,
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
	cpuInfoData, err := os.ReadFile("../../../test_data/cpuinfo/xps13-7390.txt")
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
	cpuInfoData, err := os.ReadFile("../../../test_data/cpuinfo/raspberry-pi-5.txt")
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
