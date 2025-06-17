package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/canonical/stack-utils/pkg/hardware_info"
)

func main() {
	var prettyOutput bool
	var friendlyNames bool
	var fileOutput string

	flag.BoolVar(&prettyOutput, "pretty", false, "Output pretty json. Default is compact json.")
	flag.BoolVar(&friendlyNames, "friendly", false, "Include human readable names for devices.")
	flag.StringVar(&fileOutput, "file", "", "Output json to this file. Default output is to stdout.")
	flag.Parse()

	hwInfo, err := hardware_info.Get(friendlyNames)
	if err != nil {
		log.Fatal(err)
	}

	var jsonString []byte
	if prettyOutput {
		jsonString, err = json.MarshalIndent(hwInfo, "", "  ")
		if err != nil {
			log.Fatal("Failed to marshal to JSON:", err)
		}
	} else {
		jsonString, err = json.Marshal(hwInfo)
		if err != nil {
			log.Fatal("Failed to marshal to JSON:", err)
		}
	}

	if fileOutput != "" {
		err = os.WriteFile(fileOutput, jsonString, 0644)
		if err != nil {
			log.Fatal("Failed to write to file:", err)
		}
	} else {
		fmt.Println(string(jsonString))
	}
}
