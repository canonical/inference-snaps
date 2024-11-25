package gpu

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/canonical/hardware-info/pci"
)

var lspciFiles = []string{
	"../pci/test_data/dell-precision-3660-c29399.txt",
	"../pci/test_data/dell-vostro153535-c30942.txt",
	"../pci/test_data/amd-cezanne.txt",
	"../pci/test_data/hp-elitebook845-g8-notebook-pc-c30368.txt",
	"../pci/test_data/intel-arc-a580.txt",
	"../pci/test_data/radeon_hd7450+tesla_k20xm.txt",
	"../pci/test_data/matrox_g200er2.txt",
	"../pci/test_data/rpi5.txt",
	"../pci/test_data/dell_xps13_gen10.txt",
}

func TestDisplayDevices(t *testing.T) {
	for _, lsPciFile := range lspciFiles {
		t.Run(lsPciFile, func(t *testing.T) {
			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatalf(err.Error())
			}

			pciDevices, err := pci.ParseLsPci(lsPci, true)
			if err != nil {
				t.Fatalf(err.Error())
			}

			displayDevices, err := pciGpus(pciDevices)

			jsonData, err := json.MarshalIndent(displayDevices, "", "  ")
			if err != nil {
				t.Fatalf(err.Error())
			}

			t.Log(string(jsonData))
		})
	}
}
