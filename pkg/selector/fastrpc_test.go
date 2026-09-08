package selector

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/lscompute/pkg/machine"
	lsfastrpc "github.com/canonical/lscompute/pkg/machine/device/fastrpc"
)

func TestFastRPCNPUSelection(t *testing.T) {
	machineInfo := &machine.MachineInfo{
		Devices: []any{
			lsfastrpc.Device{
				Bus:    lsfastrpc.BusName,
				Domain: lsfastrpc.CDSPDomain,
			},
		},
	}

	manifest := engines.Manifest{
		Name:    "qualcomm-npu",
		Summary: "Qualcomm Dragonwing NPU",
		Vendor:  "qualcomm",
		Devices: engines.Devices{
			Anyof: []engines.Device{{
				Type: "npu",
				Bus:  "fastrpc",
			}},
		},
	}

	score, report, err := checkEngine(machineInfo, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !report.EngineCompatible() {
		t.Fatalf("expected engine to be compatible: %+v", report)
	}
	if score <= 0 {
		t.Fatalf("expected positive score, got %d", score)
	}
}

func TestFastRPCNPURejectsNonCDSPSelection(t *testing.T) {
	machineInfo := &machine.MachineInfo{
		Devices: []any{
			lsfastrpc.Device{
				Bus:    lsfastrpc.BusName,
				Domain: lsfastrpc.ADSPDomain,
			},
		},
	}
	manifest := engines.Manifest{
		Name:    "qualcomm-npu",
		Summary: "Qualcomm Dragonwing NPU",
		Vendor:  "qualcomm",
		Devices: engines.Devices{
			Anyof: []engines.Device{{
				Type: "npu",
				Bus:  "fastrpc",
			}},
		},
	}

	score, report, err := checkEngine(machineInfo, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if report.EngineCompatible() {
		t.Fatalf("expected engine to be incompatible: %+v", report)
	}
	if score != 0 {
		t.Fatalf("expected score=0, got %d", score)
	}
}
