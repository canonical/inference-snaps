package disk

import (
	"log"

	"github.com/canonical/ml-snap-utils/pkg/types"
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
			log.Printf("%s: %s", dir, err.Error())
			continue
		}
		info[dir] = dirInfo
	}

	return info, nil
}
