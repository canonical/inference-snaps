package main

import (
	"fmt"
	"os"

	"github.com/canonical/go-snapctl"
)

func get(key string) {
	value, err := snapctl.Get(key).Run()
	if err != nil {
		fmt.Printf("Error getting value of '%s': %v\n", key, err)
		os.Exit(1)
	}

	if value == "" {
		fmt.Printf("No value set for key '%s'\n", key)
		os.Exit(1)
	}
	fmt.Println(value)
}
