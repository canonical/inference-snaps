package gpu

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/hardware_info/pci"
)

var lspciFiles = []string{
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

func TestDisplayDevices(t *testing.T) {
	for _, lsPciFile := range lspciFiles {
		t.Run(lsPciFile, func(t *testing.T) {
			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatal(err)
			}

			pciDevices, err := pci.ParseLsPci(lsPci, true)
			if err != nil {
				t.Fatal(err)
			}

			displayDevices, err := pciGpus(pciDevices)

			jsonData, err := json.MarshalIndent(displayDevices, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			t.Log(string(jsonData))
		})
	}
}
