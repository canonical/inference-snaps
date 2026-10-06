package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/utils"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/host"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

const (
	// FormatPlain is a human-readable YAML output.
	FormatPlain string = "plain"
	// FormatJSON is a machine-readable, indented JSON output with kebab-cased keys.
	FormatJSON string = "json"
)

type MachineDetails struct {
	Cpus    []CpuDetails  `json:"cpus,omitempty" yaml:"cpus,omitempty"`
	Memory  MemoryDetails `json:"memory,omitempty" yaml:"memory,omitempty"`
	Disk    []DiskDetails `json:"disks,omitempty" yaml:"disks,omitempty"`
	Devices []any         `json:"devices,omitempty" yaml:"devices,omitempty"`
}

type CpuDetails struct {
	Architecture string `json:"architecture" yaml:"architecture"`

	// amd64
	ManufacturerId string   `json:"manufacturer-id,omitempty" yaml:"manufacturer-id,omitempty"`
	Flags          []string `json:"flags,omitempty" yaml:"flags,flow,omitempty"`

	// arm64
	ImplementerId utils.HexInt `json:"implementer-id,omitempty" yaml:"implementer-id,omitempty"`
	PartNumber    utils.HexInt `json:"part-number,omitempty" yaml:"part-number,omitempty"`
	Features      []string     `json:"features,omitempty" yaml:"features,flow,omitempty"`

	// riscv64
	Isa []string `json:"isa,omitempty" yaml:"isa,omitempty"`
}

type MemoryDetails struct {
	TotalRam  uint64 `json:"total-ram" yaml:"total-ram"`
	TotalSwap uint64 `json:"total-swap" yaml:"total-swap"`
}

func (m MemoryDetails) MarshalYAML() (any, error) {
	return struct {
		TotalRam  any `yaml:"total-ram"`
		TotalSwap any `yaml:"total-swap"`
	}{
		TotalRam:  utils.FmtBytesShort(m.TotalRam),
		TotalSwap: utils.FmtBytesShort(m.TotalSwap),
	}, nil
}

type DiskDetails struct {
	MountPoint *string `json:"mount-point,omitempty" yaml:"mount-point,omitempty"`
	Path       string  `json:"path" yaml:"path"`
	Total      uint64  `json:"total" yaml:"total"`
	Avail      uint64  `json:"avail" yaml:"avail"`
}

func (d DiskDetails) MarshalYAML() (any, error) {
	return struct {
		MountPoint *string `yaml:"mount-point,omitempty"`
		Path       string  `yaml:"path"`
		Total      any     `yaml:"total"`
		Avail      any     `yaml:"avail"`
	}{
		MountPoint: d.MountPoint,
		Path:       d.Path,
		Total:      utils.FmtBytesShort(d.Total),
		Avail:      utils.FmtBytesShort(d.Avail),
	}, nil
}

type PciDeviceDetails struct {
	Bus                  string                         `json:"bus" yaml:"bus"`
	Slot                 string                         `json:"slot,omitempty" yaml:"slot,omitempty"`
	BusNumber            any                            `json:"bus-number,omitempty" yaml:"bus-number,omitempty"`
	DeviceClass          utils.HexInt                   `json:"device-class,omitempty" yaml:"device-class,omitempty"`
	ProgrammingInterface uint8                          `json:"programming-interface,omitempty" yaml:"programming-interface,omitempty"`
	VendorId             utils.HexInt                   `json:"vendor-id,omitempty" yaml:"vendor-id,omitempty"`
	DeviceId             utils.HexInt                   `json:"device-id,omitempty" yaml:"device-id,omitempty"`
	SubvendorId          utils.HexInt                   `json:"subvendor-id,omitempty" yaml:"subvendor-id,omitempty"`
	SubdeviceId          utils.HexInt                   `json:"subdevice-id,omitempty" yaml:"subdevice-id,omitempty"`
	VendorName           string                         `json:"vendor-name,omitempty" yaml:"vendor-name,omitempty"`
	DeviceName           string                         `json:"device-name,omitempty" yaml:"device-name,omitempty"`
	SubvendorName        string                         `json:"subvendor-name,omitempty" yaml:"subvendor-name,omitempty"`
	SubdeviceName        string                         `json:"subdevice-name,omitempty" yaml:"subdevice-name,omitempty"`
	AdditionalProperties *PciAdditionalDeviceProperties `json:"additional-properties,omitempty" yaml:"additional-properties,omitempty"`
}

type UsbDeviceDetails struct {
	Bus                  string            `json:"bus" yaml:"bus"`
	BusNumber            any               `json:"bus-number,omitempty" yaml:"bus-number,omitempty"`
	DeviceNumber         uint8             `json:"device-number,omitempty" yaml:"device-number,omitempty"`
	VendorId             utils.HexInt      `json:"vendor-id,omitempty" yaml:"vendor-id,omitempty"`
	ProductId            utils.HexInt      `json:"product-id,omitempty" yaml:"product-id,omitempty"`
	VendorName           string            `json:"vendor-name,omitempty" yaml:"vendor-name,omitempty"`
	ProductName          string            `json:"product-name,omitempty" yaml:"product-name,omitempty"`
	AdditionalProperties map[string]string `json:"additional-properties,omitempty" yaml:"additional-properties,omitempty"`
}

type FastRPCDeviceDetails struct {
	Bus                  string            `json:"bus" yaml:"bus"`
	Domain               string            `json:"domain,omitempty" yaml:"domain,omitempty"`
	Index                int               `json:"index,omitempty" yaml:"index,omitempty"`
	Secure               bool              `json:"secure,omitempty" yaml:"secure,omitempty"`
	AdditionalProperties map[string]string `json:"additional-properties,omitempty" yaml:"additional-properties,omitempty"`
}

type ApusysDeviceDetails struct {
	Bus        string `json:"bus" yaml:"bus"`
	Type       string `json:"type" yaml:"type"`
	SocID      string `json:"soc-id,omitempty" yaml:"soc-id,omitempty"`
	VendorName string `json:"vendor-name,omitempty" yaml:"vendor-name,omitempty"`
}

type PciAdditionalDeviceProperties struct {
	Microarchitecture string `json:"microarchitecture,omitempty" yaml:"microarchitecture,omitempty"`
	Vram              uint64 `json:"vram,omitempty" yaml:"vram,omitempty"`
	ComputeCapability string `json:"compute-capability,omitempty" yaml:"compute-capability,omitempty"`
}

func (a PciAdditionalDeviceProperties) MarshalYAML() (any, error) {
	return struct {
		Microarchitecture string `yaml:"microarchitecture,omitempty"`
		Vram              any    `yaml:"vram,omitempty"`
		ComputeCapability string `yaml:"compute-capability,omitempty"`
	}{
		Microarchitecture: a.Microarchitecture,
		Vram:              utils.FmtBytesShort(a.Vram),
		ComputeCapability: a.ComputeCapability,
	}, nil
}

func NewMachineDetails(info *machine.Machine) *MachineDetails {
	if info == nil {
		return nil
	}

	v := &MachineDetails{
		Memory: MemoryDetails(info.Memory),
	}

	// Combine all devices into a single slice
	totalDevices := len(info.PCIDevices) + len(info.USBDevices) + len(info.FastRPCDevices) + len(info.APUSYSDevices)
	v.Devices = make([]any, 0, totalDevices)

	// Add PCI devices
	for _, d := range info.PCIDevices {
		var programmingInterface uint8
		if d.ProgrammingInterface != nil {
			programmingInterface = *d.ProgrammingInterface
		}
		var subvendorId, subdeviceId uint16
		if d.SubvendorId != nil {
			subvendorId = *d.SubvendorId
		}
		if d.SubdeviceId != nil {
			subdeviceId = *d.SubdeviceId
		}
		v.Devices = append(v.Devices, PciDeviceDetails{
			Bus:                  d.Bus,
			Slot:                 d.Slot,
			BusNumber:            utils.HexInt(d.BusNumber),
			DeviceClass:          utils.HexInt(d.DeviceClass),
			ProgrammingInterface: programmingInterface,
			VendorId:             utils.HexInt(d.VendorId),
			DeviceId:             utils.HexInt(d.DeviceId),
			SubvendorId:          utils.HexInt(subvendorId),
			SubdeviceId:          utils.HexInt(subdeviceId),
			VendorName:           d.VendorName,
			DeviceName:           d.DeviceName,
			SubvendorName:        d.SubvendorName,
			SubdeviceName:        d.SubdeviceName,
			AdditionalProperties: newPciAdditionalDeviceProperties(d.AdditionalProperties),
		})
	}

	// Add USB devices
	for _, d := range info.USBDevices {
		v.Devices = append(v.Devices, UsbDeviceDetails{
			Bus:                  d.Bus,
			BusNumber:            utils.HexInt(d.BusNumber),
			DeviceNumber:         d.DeviceNumber,
			VendorId:             utils.HexInt(d.VendorId),
			ProductId:            utils.HexInt(d.ProductId),
			VendorName:           d.VendorName,
			ProductName:          d.ProductName,
			AdditionalProperties: d.AdditionalProperties,
		})
	}

	// Add FastRPC devices
	for _, d := range info.FastRPCDevices {
		v.Devices = append(v.Devices, FastRPCDeviceDetails{
			Bus:                  d.Bus,
			Domain:               string(d.Domain),
			Index:                d.Index,
			Secure:               d.Secure,
			AdditionalProperties: d.AdditionalProperties,
		})
	}

	// Add APUSYS devices
	for _, d := range info.APUSYSDevices {
		v.Devices = append(v.Devices, ApusysDeviceDetails{
			Bus:        d.Bus,
			Type:       d.Type,
			SocID:      d.SocID,
			VendorName: d.VendorName,
		})
	}

	if info.CPUs != nil {
		v.Cpus = make([]CpuDetails, len(info.CPUs))
		for i, c := range info.CPUs {
			v.Cpus[i] = CpuDetails{
				Architecture:   c.Architecture,
				ManufacturerId: c.ManufacturerId,
				Flags:          c.Flags,
				ImplementerId:  utils.HexInt(c.ImplementerId),
				PartNumber:     utils.HexInt(c.PartNumber),
				Features:       c.Features,
				Isa:            c.Isa,
			}
		}
	}

	if info.Disk != nil {
		v.Disk = make([]DiskDetails, 0, len(info.Disk))
		for _, d := range info.Disk {
			v.Disk = append(v.Disk, DiskDetails{
				MountPoint: d.MountPoint,
				Path:       d.Path,
				Total:      d.Total,
				Avail:      d.Available,
			})
		}
	}

	return v
}

func newPciAdditionalDeviceProperties(props map[string]string) *PciAdditionalDeviceProperties {
	if len(props) == 0 {
		return nil
	}

	ap := &PciAdditionalDeviceProperties{
		Microarchitecture: props["microarchitecture"],
		ComputeCapability: props["compute-capability"],
	}
	if v, ok := props["vram"]; ok {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			ap.Vram = n
		}
	}
	if *ap == (PciAdditionalDeviceProperties{}) {
		return nil
	}
	return ap
}

type machineCommand struct {
	*common.Context

	// flags
	format string
}

func Machine(ctx *common.Context) *cobra.Command {
	return newMachineCmd(ctx, "machine", "")
}

// TODO: remove when we fully migrate to "machine" command
func ShowMachine(ctx *common.Context) *cobra.Command {
	return newMachineCmd(ctx, "show-machine", `use "machine" instead`)
}

func newMachineCmd(ctx *common.Context, use, deprecated string) *cobra.Command {
	var cmd machineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               use,
		Short:             "Print information about the host machine",
		Long:              "Print information about the host machine, including hardware and compute resources",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
		Deprecated:        deprecated,
	}

	// flags
	supportedFormats := []string{"json", "yaml"}
	cobraCmd.Flags().StringVar(
		&cmd.format,
		"format",
		"yaml",
		fmt.Sprintf("output format (%s)", strings.Join(supportedFormats, ", ")),
	)

	return cobraCmd
}

func (cmd *machineCommand) run(_ *cobra.Command, _ []string) error {
	info, err := cmd.fetchMachineWithSpinner()
	if err != nil {
		return err
	}
	machineDetails := NewMachineDetails(info)
	if machineDetails == nil {
		return fmt.Errorf("failed to get machine details")
	}
	return cmd.printMachine(*machineDetails)
}

func (cmd *machineCommand) printMachine(md MachineDetails) error {
	switch cmd.format {
	case "json":
		return cmd.printMachineJson(md)
	case "yaml":
		return cmd.printMachineYaml(md)
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}
}

func (cmd *machineCommand) printMachineJson(md MachineDetails) error {
	jsonString, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return fmt.Errorf("json: %s", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *machineCommand) printMachineYaml(md MachineDetails) error {
	yamlString, err := yaml.Marshal(md)
	if err != nil {
		return fmt.Errorf("yaml: %s", err)
	}
	fmt.Printf("%s", yamlString)
	return nil
}

func (cmd *machineCommand) fetchMachineWithSpinner() (*machine.Machine, error) {
	stopProgress := common.StartProgressSpinner("Gathering machine information")
	hwInfo, warnings, err := machine.Get(host.Real(), true, true)
	stopProgress()

	if len(warnings) > 0 && cmd.Verbose {
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("getting machine info: %s", err)
	}

	return hwInfo, nil
}
