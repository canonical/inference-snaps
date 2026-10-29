// Command docgen generates the CLI reference documentation from the command tree
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/canonical/inference-snaps/modelctl/v2/commands"
	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
	"github.com/spf13/cobra/doc"
)

// rootName is a placeholder for the snap instance name
const rootName = "inference-snap"

func main() {
	out := flag.String("out", "./docs/reference/cli", "output directory")
	format := flag.String("format", "markdown", "output format: markdown or man")
	flag.Parse()

	ctx := &common.Context{
		Cache:  storage.NewMockCache(),
		Config: storage.NewMockConfig(),
	}
	// Enable all optional commands so they are documented
	root, err := commands.NewRootCmd(ctx, commands.RootOptions{Name: rootName, Chat: true, WebUi: true})
	if err != nil {
		log.Fatalf("Error building command tree: %v", err)
	}
	// Avoid a timestamp footer, so output is reproducible
	root.DisableAutoGenTag = true

	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	switch *format {
	case "markdown":
		err = doc.GenMarkdownTreeCustom(root, *out, mystLabel, dirhtmlLink)
	case "man":
		err = doc.GenManTree(root, &doc.GenManHeader{Title: strings.ToUpper(rootName), Section: "1"}, *out)
	default:
		err = fmt.Errorf("unknown format %q", *format)
	}
	if err != nil {
		log.Fatalf("Error generating docs: %v", err)
	}
}

// mystLabel prepends a MyST target, so pages can be referenced with {ref}
func mystLabel(filename string) string {
	name := strings.TrimSuffix(filepath.Base(filename), ".md")
	return fmt.Sprintf("(cli-%s)=\n\n", strings.ReplaceAll(name, "_", "-"))
}

// dirhtmlLink rewrites links to match the Sphinx dirhtml builder
func dirhtmlLink(name string) string {
	return "../" + strings.TrimSuffix(name, ".md") + "/"
}
