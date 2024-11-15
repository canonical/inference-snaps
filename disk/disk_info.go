package disk

import "log"

func Info() (map[string]*DirStats, error) {
	var info = make(map[string]*DirStats)

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
