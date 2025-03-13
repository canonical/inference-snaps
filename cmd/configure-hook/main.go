package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
	slog "github.com/canonical/go-snapctl/log"
	"github.com/canonical/ml-snap-utils/pkg/types"
)

func main() {
	slog.SetComponentName("hook.configure")

	// get stack snap option - if not set error
	stackName, err := snapctl.Get("stack").Run()
	if err != nil {
		slog.Fatalf("Error getting stack from snap options: %v", err)
	}
	if stackName == "" {
		slog.Fatal("Stack snap option is empty")
	}

	// get stacks.<new-stack> snap option and use it as definition for new stack
	stackJson, err := snapctl.Get("stacks." + stackName).Run()
	if err != nil {
		slog.Fatalf("Error getting stack definition from snap options: %v", err)
	}
	var stack types.ScoredStack
	err = json.Unmarshal([]byte(stackJson), &stack)
	if err != nil {
		slog.Fatalf("Error deserializing stack definition from snap options: %v", err)
	}

	// install components
	for _, component := range stack.Components {
		err = snapctl.InstallComponents(component).Run()
		if err != nil {
			if strings.Contains(err.Error(), "cannot install components for a snap that is unknown to the store") {
				slog.Infof("Skip component installation. Install a local build: sudo snap install %s+<path to %s>", os.Getenv("SNAP_INSTANCE_NAME"), component)
			} else if strings.Contains(err.Error(), "already installed") {
				slog.Infof("Skip component installation: already installed: %s", component)
			} else {
				slog.Fatalf("Error installing component: %v", err)
			}
		}
	}

	// set snap options from stack configurations
	for confKey, confVal := range stack.Configurations {
		valJson, err := json.Marshal(confVal)
		if err != nil {
			slog.Fatalf("Error serializing configuration %s: %v - %v", confKey, confVal, err)
		}
		err = snapctl.Set(confKey, string(valJson)).String().Run() // FIXME: for now always assume string
		if err != nil {
			slog.Fatalf("can't set snap option: %v", err)
		}
	}
}
