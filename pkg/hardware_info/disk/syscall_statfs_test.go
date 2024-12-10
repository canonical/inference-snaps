package disk

import (
	"encoding/json"
	"testing"

	"github.com/canonical/ml-snap-utils/pkg/utils"
)

var testDirs = []string{
	"/",
	"/var/lib/snapd/snaps",
}

func TestDirStats(t *testing.T) {
	for _, dir := range testDirs {
		t.Run(dir, func(t *testing.T) {
			diskStats, err := dirStats(dir)
			if err != nil {
				t.Fatalf(err.Error())
			}

			t.Log("Total:", utils.FmtGigabytes(diskStats.Total))
			t.Log("Used:", utils.FmtGigabytes(diskStats.Used))
			t.Log("Avail:", utils.FmtGigabytes(diskStats.Avail))
		})
	}
}

func TestDirStatsNonExistentDir(t *testing.T) {
	_, err := dirStats("/path/that/does/not/exist")
	if err == nil {
		t.Fatalf("Non existent dir should return error")
	}
}

func TestInfo(t *testing.T) {
	diskInfo, err := Info()
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, err := json.MarshalIndent(diskInfo, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(string(jsonData))
}
