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

func installComponents(components []string) error {
	for _, component := range components {
		stopProgress := startProgressSpinner("Installing " + component + " ")
		err := snapctl.InstallComponents(component).Run()
		stopProgress()
		if err != nil {
			if strings.Contains(err.Error(), snapdUnknownSnapError) {
				return fmt.Errorf("snap not known to the store:"+
					"\nRerun this command after manually installing %q",
					component)
			} else if strings.Contains(err.Error(), snapdTimeoutError) {
				return fmt.Errorf("timed out while installing %q:"+
					"\nMonitor the installation progress with \"snap changes\""+
					"\nRerun this command once the installation is complete",
					component)
			} else if strings.Contains(err.Error(), "already installed") {
				continue
			} else {
				return fmt.Errorf("error installing %q: %s", component, err)
			}
		}
		fmt.Println("Installed " + component)
	}

	return nil
}

func startProgressSpinner(prefix string) (stop func()) {
	s := spinner.New(spinner.CharSets[9], time.Millisecond*200)
	s.Prefix = prefix
	s.Start()

	return s.Stop
}
