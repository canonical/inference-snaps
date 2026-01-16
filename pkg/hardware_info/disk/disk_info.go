package disk

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/pkg/constants"
	"github.com/canonical/inference-snaps-cli/pkg/types"
)

var directories = []string{
	constants.SnapStoragePath,
}

// Info returns the total size and available size for root and snap dirs on the host system, using the statfs syscall.
func Info() (map[string]types.DirStats, error) {
	var info = make(map[string]types.DirStats)

	for _, dir := range directories {
		dirInfo, err := statFs(dir)
		if err != nil {
			return nil, fmt.Errorf("error getting directory info: %v", err)
		}
		info[dir] = dirInfo
	}

	return info, nil
}

// InfoFromRawData returns the total size and available size of the root and snap dirs, taking a string in which represents
// the  output of the df command.
func InfoFromRawData(dfData string) (map[string]types.DirStats, error) {
	dirInfos, err := parseDf(dfData)
	if err != nil {
		return nil, fmt.Errorf("error parsing df: %v", err)
	}

	if len(dirInfos) != len(directories) {
		return nil, fmt.Errorf("df did not return info for all dirs")
	}

	var info = make(map[string]types.DirStats)
	for i, dir := range directories {
		info[dir] = dirInfos[i]
	}

	return info, nil
}
