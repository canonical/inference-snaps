package main

import (
	"fmt"

	"github.com/canonical/stack-utils/pkg/validate"
)

func validateStackManifests(manifestFiles ...string) {
	for _, manifestPath := range manifestFiles {
		err := validate.Stack(manifestPath)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf("✅ %s\n", manifestPath)
		}
	}
}
