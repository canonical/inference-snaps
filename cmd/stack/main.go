package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/canonical/go-snapctl/env"
)

var stacksDir = env.Snap() + "/stacks"

func main() {
	// stack use [--assume-yes] [--auto]
	// stack use [--assume-yes] [<stack>]
	useCmd := flag.NewFlagSet("use", flag.ExitOnError)
	useAuto := useCmd.Bool("auto", false, "Automatically select a compatible stack")
	useAssumeYes := useCmd.Bool("assume-yes", false, "Assume yes for downloading new components")

	// stack get
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	// stack set
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	// stack unset
	unsetCmd := flag.NewFlagSet("unset", flag.ExitOnError)

	// stack load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)

	// stack download
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)

	// stack validate
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)

	// stack list [--all]
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listAll := listCmd.Bool("all", false, "Also list incompatible stacks")

	// stack info <stack>
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected a subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "use":
		useCmd.Parse(os.Args[2:])

		if *useAuto {
			if len(useCmd.Args()) != 0 {
				fmt.Println("Error: cannot specify stack with --auto flag")
				os.Exit(1)
			}
			err := autoSelectStacks(*useAssumeYes)
			if err != nil {
				fmt.Println("Error: failed to automatically set used stack:", err)
				os.Exit(1)
			}
		} else {
			stack := useCmd.Args()
			if len(stack) == 1 {
				err := useStack(stack[0], *useAssumeYes)
				if err != nil {
					fmt.Println("Error: failed use stack:", err)
					os.Exit(1)
				}
			} else if len(stack) == 0 {
				fmt.Println("Error: stack name not specified")
				os.Exit(1)
			} else {
				fmt.Println("Error: too many arguments")
				os.Exit(1)
			}
		}

	case "get":
		getCmd.Parse(os.Args[2:])
		if len(getCmd.Args()) != 1 {
			fmt.Println("Error: expected one config key as input")
			os.Exit(1)
		}
		get(getCmd.Args()[0])

	case "set":
		setCmd.Parse(os.Args[2:])
		if len(setCmd.Args()) != 1 {
			fmt.Println("Error: expected one key=value pair as input")
			os.Exit(1)
		}
		set(setCmd.Args()[0])

	case "unset":
		unsetCmd.Parse(os.Args[2:])
		if len(unsetCmd.Args()) != 1 {
			fmt.Println("Error: expected one config key as input")
			os.Exit(1)
		}
		unset(unsetCmd.Args()[0])

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

	case "list":
		listCmd.Parse(os.Args[2:])
		listStacks(*listAll)

	case "info":
		infoCmd.Parse(os.Args[2:])
		if len(infoCmd.Args()) < 1 {
			fmt.Println("Error: a stack name is required")
			os.Exit(1)
		}
		if len(infoCmd.Args()) != 1 {
			fmt.Println("Error: only one stack name can be specified")
			os.Exit(1)
		}
		stackInfo(infoCmd.Args()[0])

	default:
		fmt.Println("unexpected subcommands")
		os.Exit(1)
	}

}
