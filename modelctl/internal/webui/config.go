package webui

import (
	"fmt"
	"net/url"
)

type Config struct {
	OpenAIBaseURL string   `json:"openAIBaseURL"`
	Capabilities  []string `json:"capabilities"`
	InstanceName  string   `json:"instanceName"`
	EngineName    string   `json:"engineName"`
}

func (c Config) Validate() error {

	// Validate OpenAI base URL
	if _, err := url.Parse(c.OpenAIBaseURL); err != nil {
		return fmt.Errorf("invalid OpenAI base URL: %w", err)
	}

	return nil
}
