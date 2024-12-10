package disk

import (
	"golang.org/x/sys/unix"
)

// dirStats returns a struct with the total, used, free and available bytes for a given directory.
func dirStats(path string) (*DirStats, error) {
	var dirStats DirStats

	var fs unix.Statfs_t
	err := unix.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}

	dirStats.Total = fs.Blocks * uint64(fs.Bsize)
	dirStats.Avail = fs.Bavail * uint64(fs.Bsize)
	dirStats.Free = fs.Bfree * uint64(fs.Bsize)
	dirStats.Used = dirStats.Total - dirStats.Free
	return &dirStats, nil
}
