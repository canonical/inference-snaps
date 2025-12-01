package disk

import (
	"testing"
)

func TestHostSnapDir(t *testing.T) {
	t.Skip("See https://github.com/canonical/inference-snaps-cli/pull/237#issuecomment-3595313760")

	dfData, err := hostDf("/", "/var/snap/snapd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dfData)

	dirInfos, err := parseDf(dfData)
	if err != nil {
		t.Fatalf("can't parse df output: %v", err)
	}
	t.Log(dirInfos)
}
