package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"katemoss/common"
	"katemoss/selector"
)

func main() {
	var stacksDir string
	flag.StringVar(&stacksDir, "stacks", "stacks", "Override the path to the stacks directory.")
	flag.Parse()

	// Read json piped in from the hardware-info app
	var hardwareInfo common.HwInfo

	err := json.NewDecoder(os.Stdin).Decode(&hardwareInfo)
	if err != nil {
		log.Fatal(err)
	}

	result, err := selector.FindStack(hardwareInfo, stacksDir)
	if err != nil {
		log.Fatal(err)
	}

	resultStr, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", resultStr)
}
