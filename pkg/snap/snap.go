package snap

import (
	"github.com/canonical/go-snapctl"
	"github.com/canonical/go-snapctl/env"
)

type Snap interface {
	Restart(service ...string) error
	InstanceName() string
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
