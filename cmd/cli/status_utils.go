package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/canonical/inference-snaps-cli/pkg/engines"
)

type Status struct {
	Engine    string            `json:"engine" yaml:"engine"`
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
			break
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

	endpoints, err := serverApiUrls()
	if err != nil {
		return nil, fmt.Errorf("error getting server api endpoints: %v", err)
	}
	statusStr.Endpoints = endpoints

	return &statusStr, nil
}

func serverApiUrls() (map[string]string, error) {
	err := loadEngineEnvironment()
	if err != nil {
		return nil, fmt.Errorf("error loading engine environment: %v", err)
	}

	apiBasePath, found := os.LookupEnv(envOpenAiBasePath)
	if !found {
		return nil, fmt.Errorf("%q env var is not set", envOpenAiBasePath)
	}

	httpPortMap, err := config.Get(confHttpPort)
	if err != nil {
		return nil, fmt.Errorf("error getting %q: %v", confHttpPort, err)
	}
	httpPort := httpPortMap[confHttpPort]

	openaiUrl := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%v", httpPort),
		Path:   apiBasePath,
	}

	return map[string]string{
		// TODO add additional api endpoints like openvino on http://localhost:8080/v1
		openAi: openaiUrl.String(),
	}, nil
}
