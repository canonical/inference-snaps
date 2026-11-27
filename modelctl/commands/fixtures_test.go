package commands

import (
	"github.com/canonical/inference-snaps/modelctl/v2/internal/testfixtures"
	"github.com/canonical/lscompute/pkg/machine"
)

func machineFixture(name string) (*machine.Machine, error) {
	return testfixtures.Machine(name)
}
