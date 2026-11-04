package common

import (
	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
)

type Context struct {
	EnginesDir  string
	RuntimesDir string
	ModelsDir   string
	Version     string
	Verbose     bool
	Cache       storage.Cache
	Config      storage.Config
	Snap        snap.Snap
}
