package disk

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/canonical/stack-utils/pkg/types"
)

// hostDf runs the df command on the host to obtain stats about the provided paths.
// This only works from inside a confined snap if the mount-observe interface is connected - which is super-privileged.
func hostDf(paths ...string) (string, error) {
	// LC_ALL=POSIX df -P --block-size=1 / /var/lib/snapd/snaps
	command := exec.Command("df")
	command.Args = append(command.Args, "--portability", "--block-size=1")
	command.Args = append(command.Args, paths...)
	command.Env = append(os.Environ(), "LC_ALL=POSIX")
	out, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("df command failed: %v", err)
	}
	return string(out), nil
}

func parseDf(dfData string) ([]types.DirStats, error) {
	var parsedDirStats []types.DirStats

	lines := strings.Split(dfData, "\n")

	// Skip header line
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) != 6 {
			return nil, fmt.Errorf("not 6 columns")
		}

		totalSize, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total blocks' field: %v", err)
		}
		//usedSize, err := strconv.ParseUint(fields[2], 10, 64)
		availableSize, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'available blocks' field: %v", err)
		}

		var thisDir = types.DirStats{
			Total: totalSize,
			Avail: availableSize,
		}
		parsedDirStats = append(parsedDirStats, thisDir)
	}

	return parsedDirStats, nil
}
