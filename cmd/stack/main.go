package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl/env"
)

var stacksDir = env.Snap + "/stacks"

func main() {
	// stack select [--auto]
	// stack select [<stack>]
	selectCmd := flag.NewFlagSet("select", flag.ExitOnError)
	selectAuto := selectCmd.Bool("auto", false, "Automatically select a compatible stack")

	// stack load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)

	// stack download
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)

	// stack validate
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected a subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "select":
		selectCmd.Parse(os.Args[2:])

		if *selectAuto {
			if len(selectCmd.Args()) != 0 {
				fmt.Println("Error: cannot specify stack with --auto flag")
				os.Exit(1)
			}
			autoSelectStacks()
		} else {
			stack := selectCmd.Args()
			if len(stack) == 1 {
				selectStack(stack[0])
			} else if len(stack) == 0 {
				fmt.Println("Error: stack name not specified")
				os.Exit(1)
			} else {
				fmt.Println("Error: too many arguments")
				os.Exit(1)
			}
		}

	case "load":
		loadCmd.Parse(os.Args[2:])
		loadStacksToSnapOptions()

	case "download":
		downloadCmd.Parse(os.Args[2:])
		downloadRequiredComponents()

	// stack validate stacks/*/stack.yaml
	case "validate":
		validateCmd.Parse(os.Args[2:])
		stackFiles := validateCmd.Args()
		if len(stackFiles) == 0 {
			fmt.Println("Error: no stack manifest specified")
			os.Exit(1)
		}

		validateStackManifests(stackFiles...)

	default:
		fmt.Println("unexpected subcommands")
		os.Exit(1)
	}

}
