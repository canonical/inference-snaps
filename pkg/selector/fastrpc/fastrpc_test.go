package fastrpc

import (
	"testing"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/lscompute/pkg/machine"
	lsfastrpc "github.com/canonical/lscompute/pkg/machine/device/fastrpc"
)

func TestMatch(t *testing.T) {
	t.Run("FastRPC NPU matches CDSP", func(t *testing.T) {
		machineInfo := &machine.MachineInfo{
			Devices: []any{
				lsfastrpc.Device{
					Bus:    lsfastrpc.BusName,
					Domain: lsfastrpc.CDSPDomain,
				},
			},
		}
		device := engines.Device{Type: "npu", Bus: "fastrpc"}

		score, issues := Match(device, machineInfo)
		if score == 0 || len(issues) > 0 {
			t.Fatalf("expected a positive score with no issues, score=%d issues=%v", score, issues)
		}
	})

	t.Run("FastRPC NPU rejects non-CDSP domains", func(t *testing.T) {
		domains := []lsfastrpc.FastRPCDomain{
			lsfastrpc.ADSPDomain,
			lsfastrpc.MDSPDomain,
			lsfastrpc.SDSPDomain,
			lsfastrpc.GDSPDomain,
		}
		for _, domain := range domains {
			t.Run(string(domain), func(t *testing.T) {
				machineInfo := &machine.MachineInfo{
					Devices: []any{
						lsfastrpc.Device{
							Bus:    lsfastrpc.BusName,
							Domain: domain,
						},
					},
				}
				device := engines.Device{Type: "npu", Bus: "fastrpc"}

				score, issues := Match(device, machineInfo)
				if score != 0 {
					t.Fatalf("expected score=0, got %d", score)
				}
				if len(issues) == 0 {
					t.Fatal("expected compatibility issues")
				}
			})
		}
	})

	t.Run("FastRPC without device type accepts any domain", func(t *testing.T) {
		domains := []lsfastrpc.FastRPCDomain{
			lsfastrpc.ADSPDomain,
			lsfastrpc.MDSPDomain,
			lsfastrpc.SDSPDomain,
			lsfastrpc.CDSPDomain,
			lsfastrpc.GDSPDomain,
		}
		for _, domain := range domains {
			t.Run(string(domain), func(t *testing.T) {
				machineInfo := &machine.MachineInfo{
					Devices: []any{
						lsfastrpc.Device{
							Bus:    lsfastrpc.BusName,
							Domain: domain,
						},
					},
				}
				device := engines.Device{Bus: "fastrpc"}

				score, issues := Match(device, machineInfo)
				if score == 0 || len(issues) > 0 {
					t.Fatalf("expected a positive score with no issues, score=%d issues=%v", score, issues)
				}
			})
		}
	})

	t.Run("no FastRPC devices", func(t *testing.T) {
		machineInfo := &machine.MachineInfo{}
		device := engines.Device{Type: "npu", Bus: "fastrpc"}

		score, issues := Match(device, machineInfo)
		if score != 0 {
			t.Fatalf("expected score=0, got %d", score)
		}
		if len(issues) == 0 {
			t.Fatal("expected compatibility issues")
		}
	})

	t.Run("nil machine info", func(t *testing.T) {
		device := engines.Device{Type: "npu", Bus: "fastrpc"}

		score, issues := Match(device, nil)
		if score != 0 {
			t.Fatalf("expected score=0, got %d", score)
		}
		if len(issues) == 0 {
			t.Fatal("expected compatibility issues")
		}
	})
}
