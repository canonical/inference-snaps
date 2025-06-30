package memory

import (
	"log"
	"testing"
)

func TestInfoFromData(t *testing.T) {
	hostData, err := hostProcMemInfo()
	if err != nil {
		t.Fatalf("error getting host proc info: %v", err)
	}
	info, err := parseProcMemInfo(hostData)
	if err != nil {
		t.Fatalf("error parsing host proc info: %v", err)
	}
	log.Printf("%+v", info)
}
