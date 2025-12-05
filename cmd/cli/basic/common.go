package basic

import (
	"fmt"
	"net/url"
	"os"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/spf13/cobra"
)

const (
	groupID = "basic"

	openAi            = "openai"
	confHttpPort      = "http.port"
	envOpenAiBasePath = "OPENAI_BASE_PATH"
)

func Group(title string) *cobra.Group {
	return &cobra.Group{
		ID:    groupID,
		Title: title,
	}
}

func serverApiUrls(ctx *common.Context) (map[string]string, error) {
	err := common.LoadEngineEnvironment(ctx)
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
		openAi: openaiUrl.String(),
	}, nil
}
