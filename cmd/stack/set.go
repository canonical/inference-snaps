package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/canonical/go-snapctl"
)

func set(keyValue string) {

	if keyValue[0] == '=' {
		fmt.Println("Error: key cannot start with an equal sign")
		os.Exit(1)
	}

	// The value itself can contain an equal sign, so we split only on the first occurrence
	parts := strings.SplitN(keyValue, "=", 2)

	err := snapctl.Set(parts[0], parts[1]).Run()
	if err != nil {
		fmt.Printf("Error setting value: %v\n", err)
		os.Exit(1)
	}
}
