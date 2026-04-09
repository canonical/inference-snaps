package ui

import (
	"fmt"
	"net/url"
	"slices"
)

type Config struct {
	OpenAIBaseURL string   `json:"openAIBaseURL"`
	Capabilities  []string `json:"capabilities"`
	InstanceName  string   `json:"instanceName"`
	EngineName    string   `json:"engineName"`
}

const (
	capabilityText   string = "text"
	capabilityVision string = "vision"
)

func SupportedCapabilities() []string {
	return []string{capabilityText, capabilityVision}
}

func (c Config) Validate() error {

	// Validate OpenAI base URL
	if _, err := url.Parse(c.OpenAIBaseURL); err != nil {
		return fmt.Errorf("invalid OpenAI base URL: %w", err)
	}

	// Validate capabilities
	for _, cap := range c.Capabilities {
		if !slices.Contains(SupportedCapabilities(), cap) {
			return fmt.Errorf("unsupported capability: %q", cap)
		}
	}

	return nil
}
