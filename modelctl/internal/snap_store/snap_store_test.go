package snap_store

import "testing"

func TestSnapInfo(t *testing.T) {
	info, err := snapInfo("deepseek-r1")
	if err != nil {
		t.Fatalf("error fetching info: %v", err)
	}
	t.Logf("info: %+v", info)
}

func TestGetComponents(t *testing.T) {
	components, err := snapComponents("nr1Yeg25CSgcFmDpHD448ngXkZwSXPFA", 53, "amd64")
	if err != nil {
		t.Fatalf("error fetching components: %v", err)
	}
	t.Logf("components: %+v", components)
}
