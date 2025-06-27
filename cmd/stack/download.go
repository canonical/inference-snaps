package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/types"
)

const (
	snapdUnknownSnapError = "cannot install components for a snap that is unknown to the store"
	snapdTimeoutError     = "timeout exceeded while waiting for response"
)

func downloadRequiredComponents() {
	// get stack snap option
	stackName, err := snapctl.Get("stack").Run()
	if err != nil {
		fmt.Println("Error getting stack from snap options:", err)
		os.Exit(1)
	}
	if stackName == "" {
		fmt.Println("Stack snap option is empty")
		os.Exit(1)
	}

	// get stacks.<new-stack> snap option for the list of components
	stackJson, err := snapctl.Get("stacks." + stackName).Run()
	if err != nil {
		fmt.Println("Error getting stack definition from snap options:", err)
		os.Exit(1)
	}
	var stack types.ScoredStack
	err = json.Unmarshal([]byte(stackJson), &stack)
	if err != nil {
		fmt.Println("Error deserializing stack definition from snap options:", err)
		os.Exit(1)
	}

	// install components
	// Messages presented to the user should use the term "download" for snapctl install +component.
	for _, component := range stack.Components {
		stopProgress := startProgressDots("Downloading " + component + " ")
		err = snapctl.InstallComponents(component).Run()
		stopProgress()
		if err != nil {
			if strings.Contains(err.Error(), snapdUnknownSnapError) {
				fmt.Printf("Error: snap not known to the store. Install a local build of component: %s\n", component)
				continue
			} else if strings.Contains(err.Error(), snapdTimeoutError) {
				fmt.Printf("Error: timeout exceeded while waiting for download of %s\n", component)
				fmt.Println("Please monitor the progress using the 'snap changes' command and continue when the component installation is complete.")
				os.Exit(1)
			} else if strings.Contains(err.Error(), "already installed") {
				continue
			} else {
				fmt.Printf("Error downloading component: %s: %s\n", component, err)
				os.Exit(1)
			}
		}
		fmt.Println("Downloaded " + component)
	}
}

func startProgressDots(prefix string) (stop func()) {
	dots := []string{".", "..", "..."}
	s := spinner.New(dots, time.Second)
	s.Prefix = prefix
	s.Start()

	return s.Stop
}
