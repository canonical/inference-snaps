package types

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

var hwInfoFiles = []string{
	//"../../test_data/devices/ampere-one-m-banshee-12/hardware-info.json",
	//"../../test_data/devices/ampere-one-siryn/hardware-info.json",
	//"../../test_data/devices/ampere-one-x-banshee-8/hardware-info.json",
	"../../test_data/devices/hp-proliant-rl300-gen11-altra/hardware-info.json",
	"../../test_data/devices/hp-proliant-rl300-gen11-altra-max/hardware-info.json",
	"../../test_data/devices/i7-2600k+arc-a580/hardware-info.json",
	"../../test_data/devices/i7-10510U/hardware-info.json",
	"../../test_data/devices/mustang/hardware-info.json",
	//"../../test_data/devices/orange-pi-rv2/hardware-info.json",
	"../../test_data/devices/raspberry-pi-5/hardware-info.json",
	"../../test_data/devices/raspberry-pi-5+hailo-8/hardware-info.json",
	"../../test_data/devices/xps13-7390/hardware-info.json",
	//"../../test_data/devices/xps13-9350/hardware-info.json",
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
