package disk

import (
	"fmt"
	"os"

	"github.com/canonical/stack-utils/pkg/types"
)

func Info() (map[string]*types.DirStats, error) {
	var info = make(map[string]*types.DirStats)

	directories := []string{
		"/",
		"/var/lib/snapd/snaps", // https://snapcraft.io/docs/system-snap-directory
	}

	for _, dir := range directories {
		dirInfo, err := dirStats(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting directory stats for %s: %s\n", dir, err)
			continue
		}
		info[dir] = dirInfo
	}

	return info, nil
}
