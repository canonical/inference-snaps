package commands

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/canonical/go-snapctl/env"
	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type versionCommand struct {
	*common.Context

	// flags
	format string
}

type versionModel struct {
	Snap string `json:"snap"`
	Cli  string `json:"cli"`
}

func Version(ctx *common.Context) *cobra.Command {
	var cmd versionCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "version",
		Short:             "Show version information",
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

func (cmd *versionCommand) run(_ *cobra.Command, _ []string) error {
	versionData := cmd.getVersionData()

	switch cmd.format {
	case "json":
		jsonString, err := json.MarshalIndent(versionData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %s", err)
		}
		fmt.Printf("%s\n", jsonString)
	case "yaml":
		yamlString, err := yaml.Marshal(versionData)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %s", err)
		}
		fmt.Printf("%s", yamlString)
	default:
		return fmt.Errorf("unknown format %q", cmd.format)
	}

	return nil
}

func (cmd *versionCommand) getVersionData() versionModel {
	return versionModel{
		Snap: cleanVersionString(env.SnapVersion()),
		Cli:  cleanVersionString(getCliVersion()),
	}
}

func getCliVersion() string {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return buildInfo.Main.Version
	}
	return ""
}

func cleanVersionString(version string) string {
	if version == "" {
		return "unset"
	}
	return version
}
