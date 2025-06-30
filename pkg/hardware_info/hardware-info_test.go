package hardware_info

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/canonical/stack-utils/pkg/types"

	"github.com/go-test/deep"
)

var devices = []string{
	"hp-proliant-rl300-gen11-altra",
	"hp-proliant-rl300-gen11-altra-max",
	"i7-2600k+arc-a580",
	"mustang",
	"raspberry-pi-5",
	"raspberry-pi-5+hailo-8",
	"xps13-7390",
}

func TestGetFromFiles(t *testing.T) {
	for _, device := range devices {
		t.Run(device, func(t *testing.T) {
			hwInfo, err := GetFromRawData(t, device, true)
			if err != nil {
				t.Error(err)
			}

			var hardwareInfo types.HwInfo
			devicePath := "../../test_data/devices/" + device + "/"
			hardwareInfoData, err := os.ReadFile(devicePath + "hardware-info.json")
			if err != nil {
				t.Fatal(err)
			}
			err = json.Unmarshal(hardwareInfoData, &hardwareInfo)
			if err != nil {
				t.Fatal(err)
			}

			// Ignore friendly names during deep equal, as it depends on the version of the pci-id database
			for i := range hwInfo.PciDevices {
				hwInfo.PciDevices[i].VendorName = nil
				hwInfo.PciDevices[i].DeviceName = nil
				hwInfo.PciDevices[i].SubvendorName = nil
				hwInfo.PciDevices[i].SubdeviceName = nil
			}
			for i := range hardwareInfo.PciDevices {
				hardwareInfo.PciDevices[i].VendorName = nil
				hardwareInfo.PciDevices[i].DeviceName = nil
				hardwareInfo.PciDevices[i].SubvendorName = nil
				hardwareInfo.PciDevices[i].SubdeviceName = nil
			}

			if diff := deep.Equal(hwInfo, hardwareInfo); diff != nil {
				t.Error(diff)
			}
		})
	}
}
