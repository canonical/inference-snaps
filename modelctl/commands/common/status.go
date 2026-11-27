package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/canonical/inference-snaps/modelctl/v2/internal/engines"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/selector"
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
	for name, entrypoint := range entrypoints {
		entrypoint.Url, err = ReportedURL(entrypoint.Url)
		if err != nil {
			return nil, fmt.Errorf("formatting server entrypoint %q: %w", name, err)
		}
		entrypoints[name] = entrypoint
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
			statusStr.Notices = append(statusStr.Notices, SuggestNoActiveModel(activeEngineName))
		}
	} else if err != nil {
		return nil, fmt.Errorf("getting model status: %v", err)
	} else {
		statusStr.Model = &modelStatus
	}

	modelErr := err
	observable, err := ctx.Snap.HardwareObservable()
	if err != nil {
		return nil, fmt.Errorf("checking hardware observability: %w", err)
	}
	if !observable {
		return &statusStr, modelErr
	}

	scored, _, warnings, err := ScoreEngines(ctx)
	if ctx.Verbose {
		for _, warning := range warnings {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("scoring engines: %w", err)
	}
	if err := addBetterEngineNotice(&statusStr, scored); err != nil {
		return nil, err
	}

	return &statusStr, modelErr
}

func addBetterEngineNotice(status *Status, scored []engines.ScoredManifest) error {
	var active *engines.ScoredManifest
	for i := range scored {
		if scored[i].Name == status.Engine {
			active = &scored[i]
			break
		}
	}
	if active == nil {
		return fmt.Errorf("active engine %q not found among scored engines", status.Engine)
	}

	best, err := selector.TopEngine(scored)
	if errors.Is(err, selector.ErrorNoCompatibleEngine) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("selecting best engine: %w", err)
	}
	if best.Score > active.Score {
		status.Notices = append(status.Notices, SuggestBetterEngine())
	}
	return nil
}
