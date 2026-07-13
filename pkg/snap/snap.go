package snap

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
)

type Snap interface {
	Restart(service ...string) error
	InstanceName() string
	HardwareObservable() (bool, error)
	InstallComponent(name string) error
	ServiceStatuses() (map[string]string, error)
}

func New() Snap {
	return &snap{}
}

type snap struct{}

// Restart restarts all or a subset of snap services.
// To restart all, run without arguments.
func (*snap) Restart(service ...string) error {
	if len(service) == 0 {
		return snapctl.Restart(env.SnapName()).Run()
	}
	return snapctl.Restart(service...).Run()
}

// InstanceName returns the snap instance name.
func (*snap) InstanceName() string {
	return env.SnapInstanceName()
}

func (*snap) HardwareObservable() (bool, error) {
	connected, err := snapctl.IsConnected("hardware-observe").Run()
	if err != nil {
		return false, fmt.Errorf("checking hardware-observe connection: %w", err)
	}
	if connected {
		return true, nil
	}

	// Hardware access is available also if the snap is installed without confinement.
	// Verify by running lscpu from the system path (not staged in the snap)
	if err := exec.Command("/usr/bin/lscpu", "-V").Run(); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return false, nil
		}
		return false, fmt.Errorf("checking lscpu availability: %w", err)
	} else {
		return true, nil
	}
}

// InstallComponent installs a single snap component.
func (*snap) InstallComponent(name string) error {
	return snapctl.InstallComponents(name).Run()
}

// ServiceStatuses returns the current status of all snap services, keyed by service app name.
func (*snap) ServiceStatuses() (map[string]string, error) {
	services, err := snapctl.Services().Run()
	if err != nil {
		return nil, fmt.Errorf("getting list of services: %v", err)
	}
	statuses := make(map[string]string)
	for name, service := range services {
		// The service name is in the format <snap-name>.<service-app>, we only want the service-app part.
		_, serviceApp, found := strings.Cut(name, ".")
		if !found {
			return nil, fmt.Errorf("unexpected service name format: %q", name)
		}
		statuses[serviceApp] = service.Current
	}
	return statuses, nil
}
