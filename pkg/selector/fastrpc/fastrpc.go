package fastrpc

import (
	"fmt"
	"os"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/selector/weights"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/lscompute/pkg/machine"
	lsfastrpc "github.com/canonical/lscompute/pkg/machine/device/fastrpc"
)

func Match(manifestDevice engines.Device, machineInfo *machine.MachineInfo) (int, []string) {
	if machineInfo == nil {
		return 0, []string{"no machine info provided"}
	}

	for _, device := range machineInfo.Devices {
		fastRPCDevice, ok := device.(lsfastrpc.Device)
		if !ok || fastRPCDevice.Bus != lsfastrpc.BusName {
			continue
		}
		if manifestDevice.Type != "" && manifestDevice.Type != "npu" {
			continue
		}
		// FastRPC's CDSP is the NPU domain. A typeless FastRPC requirement is
		// intentionally broader and can match any DSP domain.
		if manifestDevice.Type == "npu" && fastRPCDevice.Domain != lsfastrpc.CDSPDomain {
			continue
		}

		for _, connection := range manifestDevice.SnapConnections {
			connected, err := checkSnapConnection(connection)
			if err != nil {
				return 0, []string{fmt.Sprintf("checking snap connection %q: %v", connection, err)}
			}
			if !connected {
				return 0, []string{fmt.Sprintf("%q is not connected", connection)}
			}
		}

		return weights.FastRPCDevice + weights.FastRPCDeviceType, nil
	}

	if manifestDevice.Type == "npu" {
		return 0, []string{"no fastrpc devices in \"cdsp\" domain on host system"}
	}
	return 0, []string{"no fastrpc devices on host system"}
}

func checkSnapConnection(connection string) (bool, error) {
	if os.Getenv("SNAP") == "" {
		// Snap connections can only be queried from inside a snap.
		return true, nil
	}
	return snap.IsConnected(connection)
}
