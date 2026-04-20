package common

import (
	"github.com/canonical/inference-snaps-cli/pkg/snap"
	"github.com/canonical/inference-snaps-cli/pkg/storage"
)

type Context struct {
	EnginesDir string
	Verbose    bool
	Cache      storage.Cache
	Config     storage.Config
	Snap       snap.Snap
}
