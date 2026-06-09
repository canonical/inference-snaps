package snap

import (
	"errors"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
)

type Snap interface {
	Restart(service ...string) error
	InstanceName() string
	HardwareObservable() (bool, error)
	InstallComponent(name string) error
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

	_, err = os.ReadDir("/sys/bus/pci/devices")
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrPermission) || os.IsPermission(err) {
		return false, nil
	}

	return false, fmt.Errorf("reading /sys/bus/pci/devices: %w", err)
}

// InstallComponent installs a single snap component.
func (*snap) InstallComponent(name string) error {
	return snapctl.InstallComponents(name).Run()
}
