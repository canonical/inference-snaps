package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/canonical/famous-models-cli/pkg/engines"
)

const (
	openAi = "openai"
)

type Status struct {
	Engine    string            `json:"engine" yaml:"engine"`
	Status    string            `json:"status" yaml:"status"`
	Endpoints map[string]string `json:"endpoints" yaml:"endpoints"`
}

func activeEngine() (*engines.ScoredManifest, error) {
	activeEngineName, err := cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("error looking up active engine: %v", err)
	}

	scoredEngines, err := scoreEngines()
	if err != nil {
		return nil, fmt.Errorf("error scoring engines: %v", err)
	}

	var scoredManifest engines.ScoredManifest
	for i := range scoredEngines {
		if scoredEngines[i].Name == activeEngineName {
			scoredManifest = scoredEngines[i]
		}
	}

	return &scoredManifest, nil
}

func statusStruct() (*Status, error) {
	var statusStr Status

	// Find the selected engine
	engine, err := activeEngine()
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

func serverApiUrls(engine *engines.ScoredManifest) (map[string]string, error) {
	// Build API URL
	apiBasePath := "v1"
	if val, ok := engine.Configurations["http.base-path"]; ok {
		apiBasePath, ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for base path: %v", val)
		}

	}
	httpPortMap, err := config.Get("http.port")
	if err != nil {
		return nil, fmt.Errorf("error getting http port: %v", err)
	}
	httpPort := httpPortMap["http.port"]
	httpPortStr := ""

	switch v := httpPort.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		httpPortStr = fmt.Sprintf("%d", v)
	case float32, float64:
		httpPortStr = fmt.Sprintf("%.0f", v)
	case string:
		httpPortStr = v
	default:
		return nil, fmt.Errorf("unexpected type for http port: %v", v)
	}

	openaiHost := fmt.Sprintf("localhost:%s", httpPortStr)
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
