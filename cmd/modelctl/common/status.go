package common

import "fmt"

type Status struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Services  map[string]string `json:"services" yaml:"services"`
	Endpoints map[string]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Model     map[string]string `json:"model,omitempty" yaml:"model,omitempty"`
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

	endpoints, err := ServerEndpoints(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	modelStatus, err := ModelStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting model status: %v", err)
	}
	statusStr.Model = modelStatus

	return &statusStr, nil
}
