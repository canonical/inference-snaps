package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// PromptYN prompts the user and returns true for 'y', false for 'n'.
func PromptYN(prompt string, defaultResponse bool) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultResponse {
			fmt.Printf("%s [Y/n] ", prompt)
		} else {
			fmt.Printf("%s [y/N] ", prompt)
		}

		input, err := reader.ReadString('\n')
		inputEnded := errors.Is(err, io.EOF)
		if err != nil && !inputEnded {
			fmt.Println()
			return false, fmt.Errorf("%w: %v", ErrNoPromptResponse, err)
		}

		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		case "":
			if !inputEnded {
				return defaultResponse, nil
			}
		default:
			if !inputEnded {
				fmt.Println(`Invalid input. Please enter "y" or "n".`)
				continue
			}
		}

		fmt.Println()
		return false, ErrNoPromptResponse
	}
}

// PromptlnEnter prompts the user for Enter in a new line
func PromptlnEnter(action string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Press [Enter] to %s, or [Ctrl+C] to abort. ", action)

	_, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nError reading input: %v\n", err)
		return false
	}

	return true
}

func PromptRestartToApplyChanges(ctx *Context, assumeYes bool) error {
	if !assumeYes {
		msg := fmt.Sprintf("Restart %s to apply the changes?", ctx.Snap.InstanceName())
		restart, err := PromptYN(msg, true)
		if err != nil {
			return fmt.Errorf("changes are saved but %s was not restarted to apply them: %w", ctx.Snap.InstanceName(), err)
		}
		if !restart {
			return nil
		}
	}

	if err := ctx.Snap.Restart(); err != nil {
		return fmt.Errorf("restarting snap: %v", err)
	}
	return nil
}
