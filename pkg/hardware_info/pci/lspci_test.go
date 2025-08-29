package pci

import (
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/utils"
)

func TestParseLsCpu(t *testing.T) {
	devices, err := utils.SubDirectories("../../../test_data/devices")
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		lsPciFile := "../../../test_data/devices/" + device + "/lspci.txt"
		t.Run(device, func(t *testing.T) {
			_, err := os.Stat(lsPciFile)
			if err != nil {
				if os.IsNotExist(err) {
					// Device does not have lspci test data, skipping
					return
				} else {
					t.Fatal(err)
				}
			}

			lsPci, err := os.ReadFile(lsPciFile)
			if err != nil {
				t.Fatal(err)
			}

			_, err = ParseLsPci(string(lsPci), true)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
