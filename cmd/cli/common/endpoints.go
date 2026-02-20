package common

import (
	"fmt"
	"net/url"
	"os"
)

const OpenAiEndpointKey = "openai"

func ServerApiUrls(ctx *Context) (map[string]string, error) {
	const (
		confHttpPort      = "http.port"
		envOpenAiBasePath = "OPENAI_BASE_PATH"
	)

	err := LoadEngineEnvironment(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading engine environment: %v", err)
	}

	apiBasePath, found := os.LookupEnv(envOpenAiBasePath)
	if !found {
		return nil, fmt.Errorf("%q env var is not set", envOpenAiBasePath)
	}

	httpPortMap, err := ctx.Config.Get(confHttpPort)
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
		OpenAiEndpointKey: openaiUrl.String(),
	}, nil
}
