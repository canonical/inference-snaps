package pci

import (
	"encoding/json"
	"os"
	"testing"
)

var testFiles = []string{
	"../../../test_data/lspci/dell-precision-3660-c29399.txt",
	"../../../test_data/lspci/dell-vostro153535-c30942.txt",
	"../../../test_data/lspci/amd-cezanne.txt",
	"../../../test_data/lspci/hp-elitebook845-g8-notebook-pc-c30368.txt",
	"../../../test_data/lspci/intel-arc-a580.txt",
	"../../../test_data/lspci/radeon_hd7450+tesla_k20xm.txt",
	"../../../test_data/lspci/matrox_g200er2.txt",
	"../../../test_data/lspci/rpi5.txt",
	"../../../test_data/lspci/dell_xps13_gen10.txt",
}

func TestParseLsCpu(t *testing.T) {
	for _, lsPciFile := range testFiles {
		t.Run(lsPciFile, func(t *testing.T) {
			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatal(err)
			}

			pciDevices, err := ParseLsPci(lsPci, true)
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
