package common

import (
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

type Context struct {
	EnginesDir  string
	RuntimesDir string
	ModelsDir   string
	Verbose     bool
	Cache       storage.Cache
	Config      storage.Config
	Snap        snap.Snap
}
