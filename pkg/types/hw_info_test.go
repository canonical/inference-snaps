package types

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/utils"
)

func TestParseHwInfo(t *testing.T) {
	machines, err := utils.SubDirectories("../../test_data/machines")
	if err != nil {
		t.Fatal(err)
	}

	for _, machine := range machines {
		hwInfoFile := "../../test_data/machines/" + machine + "/hardware-info.json"
		t.Run(machine, func(t *testing.T) {
			_, err := os.Stat(hwInfoFile)
			if err != nil {
				if os.IsNotExist(err) {
					// Device does not have hardware-info test data, skipping
					return
				} else {
					t.Fatal(err)
				}
			}

			file, err := os.Open(hwInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			var hardwareInfo HwInfo
			err = json.Unmarshal(data, &hardwareInfo)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
