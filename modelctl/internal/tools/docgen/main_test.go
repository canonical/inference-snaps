package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFormatDescription(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"plain", "Select an engine", "Select an engine"},
		{"flag", "Use --auto to select", "Use `--auto` to select"},
		{"flag at start", "--auto selects", "`--auto` selects"},
		{"adjacent flags", "--auto --fix", "`--auto` `--fix`"},
		{"bare double dash", "separate with double dashes (--) from", "separate with double dashes (`--`) from"},
		{"already code", "Use `--auto` to select", "Use `--auto` to select"},
		{"hyphenated word", "a built-in engine", "a built-in engine"},
		{"trailing whitespace", "Text. \n", "Text."},
		{"html escaped", "Print <key>", `Print \<key>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDescription(tt.in); got != tt.want {
				t.Errorf("formatDescription(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestWritePageVersion(t *testing.T) {
	var buf bytes.Buffer
	writePage(&buf, &cobra.Command{Use: "<inference-snap>"}, "v2.1.0")

	for _, want := range []string{
		"This reference describes version `v2.1.0` of the CLI.",
		"run `<inference-snap> version`.",
	} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("page is missing %q:\n%s", want, buf.String())
		}
	}
}
