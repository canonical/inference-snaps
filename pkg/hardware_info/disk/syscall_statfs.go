package disk

import (
	"github.com/canonical/ml-snap-utils/pkg/types"
	"golang.org/x/sys/unix"
)

// dirStats returns a struct with the total, used, free and available bytes for a given directory.
func dirStats(path string) (*types.DirStats, error) {
	var dirStats types.DirStats

	var fs unix.Statfs_t
	err := unix.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}

	dirStats.Total = fs.Blocks * uint64(fs.Bsize)
	dirStats.Avail = fs.Bavail * uint64(fs.Bsize)
	return &dirStats, nil
}
