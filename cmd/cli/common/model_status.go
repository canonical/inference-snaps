package common

import (
	"fmt"
	"strings"

	"github.com/canonical/inference-snaps-cli/pkg/models"
)

func ModelStatus(ctx *Context) (map[string]string, error) {
	activeModelId, err := ctx.Cache.GetActiveModel()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveModel, err)
	}

	activeModelManifest, err := models.LoadManifest(ctx.ModelsDir, activeModelId)
	if err != nil {
		return nil, fmt.Errorf("loading model manifest: %v", err)
	}

	status := make(map[string]string)
	for _, kv := range activeModelManifest.Environment {
		// Split into key/value
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return status, fmt.Errorf("invalid env var %q", kv)
		}
		k, v := parts[0], parts[1]

		if k == "MODEL_NAME" {
			status["name"] = v
		}
	}

	return status, nil
}
