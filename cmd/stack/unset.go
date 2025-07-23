package main

import (
	"fmt"
	"os"

	"github.com/canonical/go-snapctl"
)

func unset(key string) {

	err := snapctl.Unset(key).Run()
	if err != nil {
		fmt.Printf("Error unsetting value: %v\n", err)
		os.Exit(1)
	}
}
