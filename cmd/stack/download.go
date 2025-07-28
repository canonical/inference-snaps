package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/canonical/go-snapctl"
)

const (
	snapdUnknownSnapError = "cannot install components for a snap that is unknown to the store"
	snapdTimeoutError     = "timeout exceeded while waiting for response"
)

func downloadComponents(components []string) error {
	// install components
	// Messages presented to the user should use the term "download" for snapctl install +component.
	for _, component := range components {
		stopProgress := startProgressDots("Downloading " + component + " ")
		err := snapctl.InstallComponents(component).Run()
		stopProgress()
		if err != nil {
			if strings.Contains(err.Error(), snapdUnknownSnapError) {
				fmt.Println("Error: snap not known to the store. Install a local build of component: %s", component)
				continue
			} else if strings.Contains(err.Error(), snapdTimeoutError) {
				msg := "timeout exceeded while waiting for download of: " + component +
					"\nPlease monitor the progress using the 'snap changes' command and continue when the component installation is complete."
				return fmt.Errorf(msg)
			} else if strings.Contains(err.Error(), "already installed") {
				continue
			} else {
				return fmt.Errorf("error downloading component: %s: %s", component, err)
			}
		}
		fmt.Println("Downloaded " + component)
	}

	return nil
}

func startProgressDots(prefix string) (stop func()) {
	dots := []string{".", "..", "..."}
	s := spinner.New(dots, time.Second)
	s.Prefix = prefix
	s.Start()

	return s.Stop
}
