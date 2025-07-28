package main

import (
	"flag"
	"fmt"
	"os"

	_validate "github.com/canonical/stack-utils/pkg/validate"
)

func validate(args []string) error {
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	validateCmd.Parse(os.Args[2:])
	stackFiles := validateCmd.Args()

	if len(stackFiles) == 0 {
		return fmt.Errorf("no stack manifest specified")
	}

	return validateStackManifests(stackFiles...)
}

func validateStackManifests(manifestFiles ...string) error {
	for _, manifestPath := range manifestFiles {
		err := _validate.Stack(manifestPath)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", manifestPath, err)
		} else {
			fmt.Printf("✅ %s\n", manifestPath)
		}
	}

	return nil
}
