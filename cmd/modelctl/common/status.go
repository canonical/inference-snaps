package common

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
)

type Status struct {
	Engine      string             `json:"engine" yaml:"engine"`
	Services    map[string]string  `json:"services" yaml:"services"`
	Entrypoints Entrypoints        `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Model       *map[string]string `json:"model" yaml:"model"`
	Notices     []string           `json:"notices,omitempty" yaml:"notices,omitempty"`
}

func SnapStatus(ctx *Context) (*Status, error) {
	var statusStr Status

	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", LookingUpActiveEngine, err)
	}
	if activeEngineName == "" {
		return nil, ErrNoActiveEngine
	}
	statusStr.Engine = activeEngineName

	services, err := ctx.Snap.ServiceStatuses()
	if err != nil {
		return nil, fmt.Errorf("getting service statuses: %v", err)
	}
	statusStr.Services = services

	entrypoints, err := ServerEntrypoints(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting server entrypoints: %v", err)
	}
	statusStr.Entrypoints = entrypoints

	modelStatus, err := ModelStatus(ctx)
	if err == ErrNoActiveModel {
		// ErrNoActiveModel is also the normal state for engines that define no
		// model, so only flag disk exhaustion when a model was actually expected.
		engineManifest, mErr := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
		if mErr != nil {
			return nil, fmt.Errorf("loading engine manifest: %w", mErr)
		}
		if len(engineManifest.Model.Options) > 0 {
			statusStr.Notices = append(statusStr.Notices, SuggestNotEnoughSpaceForModel(activeEngineName))
		}
		return &statusStr, err
	} else if err != nil {
		return nil, fmt.Errorf("getting model status: %v", err)
	}
	statusStr.Model = &modelStatus

	return &statusStr, nil
}
