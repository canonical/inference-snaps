package common

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	openAiEndpointKey = "openai"
	protocolKey       = "protocol"
	basePathKey       = "base-path"
)

func ServerEndpoints(ctx *Context) (map[string]string, error) {
	settings, err := EngineComponentSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading engine environment: %v", err)
	}
	return serverEndpoints(ctx, settings)
}

func serverEndpoints(ctx *Context, settingsCollection []ComponentSettings) (map[string]string, error) {
	endpoints := make(map[string]string)
	for _, settings := range settingsCollection {

		// TODO: Remove this check in a future release
		for _, env := range settings.Environment {
			if strings.HasPrefix(env, "OPENAI_BASE_PATH") {
				return nil, fmt.Errorf("OPENAI_BASE_PATH env in component %q is deprecated; set server settings in \"servers\".",
					settings.componentName)
			}
		}

		for serverName, serverSettings := range settings.Servers {
			switch serverSettings[protocolKey] {
			case "http", "https":
				httpUrl, err := serverHttpUrl(ctx, serverSettings)
				if err != nil {
					return nil, fmt.Errorf("getting server HTTP URL: %v", err)
				}
				endpoints[serverName] = httpUrl
			default:
				return nil, fmt.Errorf("unsupported protocol %q for server %q in component %q",
					serverSettings["protocol"], serverName, settings.componentName)
			}
		}

		// If builtin webui is enabled, also list it as an endpoint
		if WebUiEnabled() {
			webUiUrl, err := UiServerHttpUrl(ctx)
			if err != nil {
				return nil, fmt.Errorf("getting web UI url: %v", err)
			}
			endpoints["webui"] = webUiUrl
		}
	}

	return endpoints, nil
}

func serverHttpUrl(ctx *Context, serverConfig map[string]string) (string, error) {
	const (
		confHttpPort    = "http.port"
		defaultBasePath = "/"
	)

	httpPortMap, err := ctx.Config.Get(confHttpPort)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", confHttpPort, err)
	}
	httpPort := httpPortMap[confHttpPort]

	basePath, found := serverConfig[basePathKey]
	if !found {
		basePath = defaultBasePath
	}

	endpointUrl := url.URL{
		Scheme: serverConfig[protocolKey],
		Host:   fmt.Sprintf("localhost:%v", httpPort),
		Path:   basePath,
	}

	return endpointUrl.String(), nil
}

func OpenAiEndpoint(ctx *Context) (string, error) {
	serverEndpoints, err := ServerEndpoints(ctx)
	if err != nil {
		return "", fmt.Errorf("getting server endpoints: %v", err)
	}
	openaiEndpoint, found := serverEndpoints[openAiEndpointKey]
	if !found {
		return "", fmt.Errorf("%q not found in server endpoints", openAiEndpointKey)
	}
	return openaiEndpoint, nil
}

func UiServerHttpUrl(ctx *Context) (string, error) {
	const (
		confWebuiHttpPort = "webui.http.port"
		defaultBasePath   = "/"
	)

	httpPortMap, err := ctx.Config.Get(confWebuiHttpPort)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", confWebuiHttpPort, err)
	}
	httpPort := httpPortMap[confWebuiHttpPort]

	endpointUrl := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%v", httpPort),
		Path:   defaultBasePath,
	}

	return endpointUrl.String(), nil
}
