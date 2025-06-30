package disk

import (
	"encoding/json"
	"testing"

	"github.com/canonical/stack-utils/pkg/utils"
)

var testDirs = []string{
	"/",
	"/var/lib/snapd/snaps",
}

func TestDirStats(t *testing.T) {
	for _, dir := range testDirs {
		t.Run(dir, func(t *testing.T) {
			diskStats, err := statFs(dir)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("Total:", utils.FmtGigabytes(diskStats.Total))
			t.Log("Avail:", utils.FmtGigabytes(diskStats.Avail))
		})
	}
}

func TestDirStatsNonExistentDir(t *testing.T) {
	_, err := statFs("/path/that/does/not/exist")
	if err == nil {
		t.Fatal("Non existent dir should return error")
	}
}

func TestInfo(t *testing.T) {
	diskInfo, err := Info()
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(diskInfo, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(jsonData))
}
