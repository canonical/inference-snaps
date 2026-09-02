package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/lscompute/pkg/machine"
	"github.com/canonical/lscompute/pkg/machine/host"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

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
	info, err := cmd.fetchMachineInfoWithSpinner()
	if err != nil {
		return err
	}

	return cmd.printMachineInfo(info)
}

func (cmd *machineCommand) printMachineInfo(info *machine.MachineInfo) error {
	switch cmd.format {
	case "json":
		return cmd.printMachineInfoJson(info)
	case "yaml":
		return cmd.printMachineInfoYaml(info)
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}
}

func (cmd *machineCommand) printMachineInfoJson(info *machine.MachineInfo) error {
	jsonString, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("json: %s", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *machineCommand) printMachineInfoYaml(info *machine.MachineInfo) error {
	yamlString, err := yaml.Marshal(info)
	if err != nil {
		return fmt.Errorf("yaml: %s", err)
	}
	fmt.Printf("%s", yamlString)
	return nil
}

func (cmd *machineCommand) fetchMachineInfoWithSpinner() (*machine.MachineInfo, error) {
	stopProgress := common.StartProgressSpinner("Gathering machine information")
	hwInfo, warnings, err := machine.Get(host.Real(), true)
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
