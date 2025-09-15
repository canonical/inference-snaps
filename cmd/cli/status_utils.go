package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/canonical/go-snapctl"
	"github.com/canonical/stack-utils/pkg/engines"
)

const (
	openAi = "openai"
)

type Status struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Status    string            `json:"status" yaml:"status"`
	Endpoints map[string]string `json:"endpoints" yaml:"endpoints"`
}

func scoredEnginesFromOptions() ([]engines.ScoredManifest, error) {
	engineJson, err := snapctl.Get("engines").Document().Run()
	if err != nil {
		return nil, fmt.Errorf("error loading engines: %v", err)
	}

	engines, err := parseEnginesJson(engineJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing engines: %v", err)
	}

	return engines, nil
}

func selectedEngineFromOptions() (engines.ScoredManifest, error) {
	selectedEngineName, err := snapctl.Get("engine").Run()
	if err != nil {
		return engines.ScoredManifest{}, fmt.Errorf("error loading selected engine: %v", err)
	}

	enginesJson, err := snapctl.Get("engines." + selectedEngineName).Document().Run()
	if err != nil {
		return engines.ScoredManifest{}, fmt.Errorf("error loading engine: %v", err)
	}

	engine, err := parseEngineJson(enginesJson)
	if err != nil {
		return engines.ScoredManifest{}, fmt.Errorf("error parsing engine: %v", err)
	}

	return engine, nil
}

func statusStruct() (*Status, error) {
	var statusStr Status

	// Find the selected engine
	engine, err := selectedEngineFromOptions()
	if err != nil {
		return nil, fmt.Errorf("error loading selected engine: %v", err)
	}
	statusStr.Engine = engine.Name

	ssc, err := serverStatusCode(engine.Name)
	if err != nil {
		return nil, fmt.Errorf("error getting server status: %v", err)
	}
	switch ssc {
	case 0:
		statusStr.Status = "online"
	case 1:
		statusStr.Status = "starting"
	case 2:
		statusStr.Status = "offline"
	default:
		statusStr.Status = "unknown"
	}

	endpoints, err := serverApiUrls(engine)
	if err != nil {
		return nil, fmt.Errorf("error getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	return &statusStr, nil
}

func serverApiUrls(engine engines.ScoredManifest) (map[string]string, error) {
	// Build API URL
	apiBasePath := "v1"
	if val, ok := engine.Configurations["http.base-path"]; ok {
		apiBasePath, ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for base path: %v", val)
		}

	}
	httpPort, err := snapctl.Get("http.port").Run()
	if err != nil {
		return nil, fmt.Errorf("error getting http port: %v", err)
	}

	openaiHost := fmt.Sprintf("localhost:%s", httpPort)
	openaiUrl := url.URL{Scheme: "http", Host: openaiHost, Path: apiBasePath}
	return map[string]string{openAi: openaiUrl.String()}, nil

	// TODO add additional api endpoints like openvino on http://localhost:8080/v1
}

func serverStatusCode(engineName string) (int, error) {
	// Depend on existing check server scripts for status
	checkScript := os.ExpandEnv("$SNAP/engines/" + engineName + "/check-server")
	cmd := exec.Command(checkScript)
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("error checking server: %v", err)
	}

	checkExitCode := 0
	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			checkExitCode = exitError.ExitCode()
		}
	}
	return checkExitCode, nil
}
