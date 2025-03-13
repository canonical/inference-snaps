package memory

import (
	"encoding/json"
	"testing"
)

func TestInfo(t *testing.T) {
	info, err := Info()
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(jsonData))
}
