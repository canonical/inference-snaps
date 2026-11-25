package common

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
)

func withStdin(input string, fn func()) {
	originalStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	if _, err := w.WriteString(input); err != nil {
		panic(err)
	}
	_ = w.Close()

	os.Stdin = r
	defer func() {
		os.Stdin = originalStdin
		_ = r.Close()
	}()

	fn()
}

func printToStdout(a ...any) {
	fmt.Println(append([]any{"->"}, a...)...)
}

func ExamplePromptYN_defaultYes() {
	withStdin("\n", func() {
		printToStdout(PromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] -> true <nil>
}

func ExamplePromptYN_invalidThenNo() {
	withStdin("maybe\nn\n", func() {
		printToStdout(PromptYN("Proceed?", true))
	})

	// Output:
	// Proceed? [Y/n] Invalid input. Please enter "y" or "n".
	// Proceed? [Y/n] -> false <nil>
}

func ExamplePromptYN_finalLineWithoutNewline() {
	withStdin("y", func() {
		printToStdout(PromptYN("Proceed?", false))
	})

	// Output:
	// Proceed? [y/N] -> true <nil>
}

func TestPromptYNWithoutResponse(t *testing.T) {
	tests := map[string]struct {
		stdin           string
		defaultResponse bool
	}{
		"closed input, default yes":        {stdin: "", defaultResponse: true},
		"closed input, default no":         {stdin: "", defaultResponse: false},
		"blank final line, default yes":    {stdin: "  ", defaultResponse: true},
		"invalid final line":               {stdin: "maybe", defaultResponse: true},
		"invalid answer then closed input": {stdin: "maybe\n", defaultResponse: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			withStdin(tc.stdin, func() {
				answer, err := PromptYN("Proceed?", tc.defaultResponse)
				if !errors.Is(err, ErrNoPromptResponse) {
					t.Errorf("expected %q, got: %v", ErrNoPromptResponse, err)
				}
				if answer {
					t.Error("expected false without a response")
				}
			})
		})
	}

	t.Run("unreadable input", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		_ = w.Close()
		_ = r.Close()

		originalStdin := os.Stdin
		os.Stdin = r
		t.Cleanup(func() { os.Stdin = originalStdin })

		if _, err := PromptYN("Proceed?", true); !errors.Is(err, ErrNoPromptResponse) {
			t.Errorf("expected %q, got: %v", ErrNoPromptResponse, err)
		}
	})
}

func ExamplePromptlnEnter() {
	withStdin("\n", func() {
		printToStdout(PromptlnEnter("continue"))
	})

	// Output:
	// Press [Enter] to continue, or [Ctrl+C] to abort. -> true
}

func ExamplePromptRestartToApplyChanges_confirmed() {
	withStdin("\n", func() {
		printToStdout(PromptRestartToApplyChanges(&Context{Snap: snap.Mock()}, false))
	})

	// Output:
	// Restart mock-snap to apply the changes? [Y/n] [mock] Restarting all services
	// -> <nil>
}

func ExamplePromptRestartToApplyChanges_declined() {
	withStdin("n\n", func() {
		printToStdout(PromptRestartToApplyChanges(&Context{Snap: snap.Mock()}, false))
	})

	// Output:
	// Restart mock-snap to apply the changes? [Y/n] -> <nil>
}

func ExamplePromptRestartToApplyChanges_assumeYes() {
	withStdin("", func() {
		printToStdout(PromptRestartToApplyChanges(&Context{Snap: snap.Mock()}, true))
	})

	// Output:
	// [mock] Restarting all services
	// -> <nil>
}

func TestPromptRestartToApplyChangesWithoutResponse(t *testing.T) {
	withStdin("", func() {
		err := PromptRestartToApplyChanges(&Context{Snap: snap.Mock()}, false)
		if !errors.Is(err, ErrNoPromptResponse) {
			t.Errorf("expected %q, got: %v", ErrNoPromptResponse, err)
		}
	})
}
