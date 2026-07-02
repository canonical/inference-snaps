package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/spf13/cobra"
)

type exportStatusCommand struct {
	*common.Context
}

func ExportStatus(ctx *common.Context) *cobra.Command {
	var cmd exportStatusCommand
	cmd.Context = ctx

	cobraCmd := &cobra.Command{
		Use:               "export-status [<dir>]",
		Short:             "Export the current status",
		Long:              "Write the current status to status.json. It is stored in $SNAP_COMMON/share, unless a directory is passed as the first argument.",
		Hidden:            true,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              cmd.run,
	}

	return cobraCmd
}

type exportedStatus struct {
	Endpoints map[string]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Model     map[string]string `json:"model,omitempty" yaml:"model,omitempty"`
}

// Deprecated: exportedOpenai is to be deprecated in favor of exportedStatus.
// The Open WebUI snap first needs to be updated to consume exportedStatus
type exportedOpenai struct {
	BaseUrl string `json:"base_url"`
}

func (cmd *exportStatusCommand) run(_ *cobra.Command, args []string) error {

	statusStr, err := common.SnapStatus(cmd.Context)
	if err != nil {
		return fmt.Errorf("getting status: %v", err)
	}

	// Decouple internal status definition from shared one
	sharedStatusStr := &exportedStatus{
		Endpoints: statusStr.Endpoints,
		Model:     statusStr.Model,
	}

	var shareDir string
	if len(args) > 0 {
		shareDir = args[0]
	} else {
		shareDir = filepath.Join(os.Getenv("SNAP_COMMON"), "share")
	}

	return cmd.writeShareFiles(sharedStatusStr, shareDir)
}

// writeShareFiles writes status.json (and optionally openai.json) to shareDir.
func (cmd *exportStatusCommand) writeShareFiles(status *exportedStatus, shareDir string) error {
	statusJson, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("json: %v", err)
	}

	if err := os.MkdirAll(shareDir, 0o755); err != nil {
		return fmt.Errorf("creating status share directory: %v", err)
	}

	statusFilePath := filepath.Join(shareDir, "status.json")
	tmpStatusFilePath := statusFilePath + ".tmp"
	if err := os.WriteFile(tmpStatusFilePath, statusJson, 0o644); err != nil {
		return fmt.Errorf("writing status.json: %v", err)
	}
	if err := os.Rename(tmpStatusFilePath, statusFilePath); err != nil {
		return fmt.Errorf("renaming status.json: %v", err)
	}

	// Deprecated: Write openai.json for backwards compatibility while Open WebUI snap is not updated
	openaiFilePath := filepath.Join(shareDir, "openai.json")
	if endpoint, ok := status.Endpoints["openai"]; ok {
		myOpenaiConfig := &exportedOpenai{
			BaseUrl: endpoint,
		}
		openaiJson, err := json.Marshal(myOpenaiConfig)
		if err != nil {
			return fmt.Errorf("json: %v", err)
		}
		tmpOpenaiFilePath := openaiFilePath + ".tmp"
		if err := os.WriteFile(tmpOpenaiFilePath, openaiJson, 0o644); err != nil {
			return fmt.Errorf("writing openai.json: %v", err)
		}
		if err := os.Rename(tmpOpenaiFilePath, openaiFilePath); err != nil {
			return fmt.Errorf("renaming openai.json: %v", err)
		}
	} else if err := os.Remove(openaiFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing openai.json: %v", err)
	}

	return nil
}
