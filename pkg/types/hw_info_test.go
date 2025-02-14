package types

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

var hwInfoFiles = []string{
	"../../test_data/hardware_info/amd-ryzen7-5700g.json",
	"../../test_data/hardware_info/amd-ryzen9-7900.json",
	"../../test_data/hardware_info/dell-r730xd.json",
	"../../test_data/hardware_info/hp-dl380p-gen8.json",
	"../../test_data/hardware_info/i7-2600k.json",
	"../../test_data/hardware_info/nuc11-i5-1145G7.json",
	"../../test_data/hardware_info/xps13-gen10.json",
}

func TestParseHwInfo(t *testing.T) {
	for _, hwInfoFile := range hwInfoFiles {
		t.Run(hwInfoFile, func(t *testing.T) {
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
