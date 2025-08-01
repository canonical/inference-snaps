package disk

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/types"
	"golang.org/x/sys/unix"
)

// statFs returns a struct with the total, used, free and available bytes for a given directory.
func statFs(path string) (types.DirStats, error) {
	var pathStats types.DirStats

	var fs unix.Statfs_t
	err := unix.Statfs(path, &fs)
	if err != nil {
		return pathStats, fmt.Errorf("statfs failed: %v", err)
	}

	pathStats.Total = fs.Blocks * uint64(fs.Bsize)
	pathStats.Avail = fs.Bavail * uint64(fs.Bsize)
	return pathStats, nil
}
