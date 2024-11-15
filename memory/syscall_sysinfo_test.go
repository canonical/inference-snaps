package memory

import (
	"encoding/json"
	"testing"
)

func TestInfo(t *testing.T) {
	info, err := Info()
	if err != nil {
		t.Fatalf(err.Error())
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(string(jsonData))
}
