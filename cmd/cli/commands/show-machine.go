package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/canonical/inference-snaps-cli/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type showMachineCommand struct {
	*common.Context

	// flags
	format string
}

func ShowMachine(ctx *common.Context) *cobra.Command {
	var cmd showMachineCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "show-machine",
		Short:             "Print information about the host machine",
		Long:              "Print information about the host machine, including hardware and compute resources",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
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

func (cmd *showMachineCommand) run(_ *cobra.Command, _ []string) error {
	info, err := cmd.fetchMachineInfoWithSpinner()
	if err != nil {
		return err
	}

	return cmd.printMachineInfo(info)
}

func (cmd *showMachineCommand) printMachineInfo(info *types.HwInfo) error {
	switch cmd.format {
	case "json":
		return cmd.printMachineInfoJson(info)
	case "yaml":
		return cmd.printMachineInfoYaml(info)
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}
}

func (cmd *showMachineCommand) printMachineInfoJson(info *types.HwInfo) error {
	jsonString, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("json: %s", err)
	}
	fmt.Printf("%s\n", jsonString)
	return nil
}

func (cmd *showMachineCommand) printMachineInfoYaml(info *types.HwInfo) error {
	yamlString, err := yaml.Marshal(info)
	if err != nil {
		return fmt.Errorf("yaml: %s", err)
	}
	fmt.Printf("%s", yamlString)
	return nil
}

func (cmd *showMachineCommand) fetchMachineInfoWithSpinner() (*types.HwInfo, error) {
	stopProgress := common.StartProgressSpinner("Gathering machine information")
	hwInfo, warnings, err := hardware_info.Get(true)
	stopProgress()

	if len(warnings) > 0 {
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("getting machine info: %s", err)
	}

	return hwInfo, nil
}
