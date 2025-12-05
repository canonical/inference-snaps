package others

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/hardware_info"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type showMachineCommand struct {
	*common.Context

	// flags
	format string
}

func ShowMachineCommand(ctx *common.Context) *cobra.Command {
	var cmd showMachineCommand
	cmd.Context = ctx

	cobra := &cobra.Command{
		Use:               "show-machine",
		Short:             "Print information about the host machine",
		Long:              "Print information about the host machine, including hardware and compute resources",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	// flags
	cobra.Flags().StringVar(&cmd.format, "format", "yaml", "output format")

	return cobra
}

func (cmd *showMachineCommand) run(_ *cobra.Command, _ []string) error {
	hwInfo, err := hardware_info.Get(true)
	if err != nil {
		return fmt.Errorf("failed to get machine info: %s", err)
	}

	switch cmd.format {
	case "json":
		jsonString, err := json.MarshalIndent(hwInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml":
		yamlString, err := yaml.Marshal(hwInfo)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		fmt.Printf("%s", yamlString)
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}
