package pci

import (
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

func TestParseLsCpu(t *testing.T) {
	machines, err := utils.SubDirectories("../../../test_data/machines")
	if err != nil {
		t.Fatal(err)
	}

	for _, machine := range machines {
		lsPciFile := "../../../test_data/machines/" + machine + "/lspci.txt"
		t.Run(machine, func(t *testing.T) {
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
