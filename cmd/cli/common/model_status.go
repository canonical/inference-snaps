package common

import (
	"fmt"
	"strings"
)

func ModelStatus(ctx *Context) (map[string]string, error) {
	settings, err := EngineComponentSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading engine environment: %v", err)
	}
	return modelStatus(settings)
}

func modelStatus(settingsCollection []ComponentSettings) (map[string]string, error) {
	status := make(map[string]string)
	for _, settings := range settingsCollection {
		for i := range settings.Environment {
			// Split into key/value
			kv := settings.Environment[i]
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) != 2 {
				return status, fmt.Errorf("invalid env var %q", kv)
			}
			k, v := parts[0], parts[1]

			if k == "MODEL_NAME" {
				status["name"] = v
			}
		}
	}

	return status, nil
}
