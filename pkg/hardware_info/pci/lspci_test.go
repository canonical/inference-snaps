package pci

import (
	"encoding/json"
	"os"
	"testing"
)

var testFiles = []string{
	"../../../test_data/devices/ampere-one-m-banshee-12/lspci.txt",
	"../../../test_data/devices/ampere-one-siryn/lspci.txt",
	"../../../test_data/devices/ampere-one-x-banshee-8/lspci.txt",
	"../../../test_data/devices/hp-proliant-rl300-gen11-altra/lspci.txt",
	"../../../test_data/devices/hp-proliant-rl300-gen11-altra-max/lspci.txt",
	"../../../test_data/devices/i7-2600k+arc-a580/lspci.txt",
	"../../../test_data/devices/i7-10510U/lspci.txt",
	"../../../test_data/devices/mustang/lspci.txt",
	"../../../test_data/devices/orange-pi-rv2/lspci.txt",
	"../../../test_data/devices/raspberry-pi-5/lspci.txt",
	"../../../test_data/devices/raspberry-pi-5+hailo-8/lspci.txt",
	"../../../test_data/devices/xps13-7390/lspci.txt",
	"../../../test_data/devices/xps13-9350/lspci.txt",
}

func TestParseLsCpu(t *testing.T) {
	for _, lsPciFile := range testFiles {
		t.Run(lsPciFile, func(t *testing.T) {
			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatal(err)
			}

			pciDevices, err := ParseLsPci(string(lsPci), true)
			if err != nil {
				t.Fatal(err)
			}

			jsonData, err := json.MarshalIndent(pciDevices, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			t.Log(string(jsonData))
		})
	}
}
