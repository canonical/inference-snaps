package main

import (
	"encoding/json"
	"fmt"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/selector"
	"github.com/canonical/stack-utils/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Initialize snap configurations",
		// Long:  "",
		Hidden:            true, // command for internal use
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              load,
	}
	rootCmd.AddCommand(cmd)
}

func load(_ *cobra.Command, _ []string) error {
	if !utils.IsRootUser() {
		return ErrPermissionDenied
	}
	return loadEnginesToSnapOptions()
}

func loadEnginesToSnapOptions() error {
	fmt.Println("Loading engines to snap options ...")

	allEngines, err := selector.LoadManifestsFromDir(enginesDir)
	if err != nil {
		return fmt.Errorf("error loading engines: %v", err)
	}

	// set all engines as snap options
	// TODO: change to also handle engine deletions
	for _, engine := range allEngines {
		engineJson, err := json.Marshal(engine)
		if err != nil {
			return fmt.Errorf("error serializing engines: %s", err)
		}

		err = snapctl.Set("engines."+engine.Name, string(engineJson)).Document().Run()
		if err != nil {
			return fmt.Errorf("error setting engines option: %s", err)
		}
	}

	return nil
}
