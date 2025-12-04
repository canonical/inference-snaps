package engine

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/canonical/inference-snaps-cli/pkg/engines"
	"github.com/canonical/inference-snaps-cli/pkg/selector"
	"github.com/spf13/cobra"
)

const groupID = "engine"

func Group(title string) *cobra.Group {
	return &cobra.Group{
		ID:    groupID,
		Title: title,
	}
}

func scoreEngines(ctx *common.Context) ([]engines.ScoredManifest, error) {
	allEngines, err := selector.LoadManifestsFromDir(ctx.EnginesDir)
	if err != nil {
		return nil, fmt.Errorf("error loading engines: %v", err)
	}

	machineInfo, err := ctx.Cache.GetMachineInfo()
	if err != nil {
		return nil, fmt.Errorf("error getting machine info: %v", err)
	}

	// score engines
	scoredEngines, err := selector.ScoreEngines(machineInfo, allEngines)
	if err != nil {
		return nil, fmt.Errorf("error scoring engines: %v", err)
	}

	return scoredEngines, nil
}
