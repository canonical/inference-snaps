package common

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
)

const (
	openAiEndpointKey = "openai"
)

func ServerEndpoints(ctx *Context) (map[string]string, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", LookingUpActiveEngine, err)
	}
	activeEngineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading active engine manifest: %v", err)
	}

	// If the engine does not list a runtime, return no endpoints
	if activeEngineManifest.Runtime == "" {
		return nil, nil
	}

	runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, activeEngineManifest.Runtime)
	if err != nil {
		return nil, fmt.Errorf("loading runtime manifest: %v", err)
	}

	endpoints := make(map[string]string)

	for serverName, serverSettings := range runtimeManifest.Servers {
		switch serverSettings.Protocol {
		case "http", "https":
			httpUrl, err := serverHttpUrl(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("getting server HTTP URL: %v", err)
			}
			endpoints[serverName] = httpUrl
		default:
			return nil, fmt.Errorf("unsupported protocol %q for server %q in component %q",
				serverSettings.Protocol, serverName, activeEngineManifest.Runtime)
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

	return endpoints, nil
}

func serverHttpUrl(ctx *Context, server runtimes.Server) (string, error) {
	const (
		confHttpPort    = "http.port"
		defaultBasePath = "/"
		confHost        = "http.host"
	)

	httpPortMap, err := ctx.Config.Get(confHttpPort)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", confHttpPort, err)
	}
	httpPort := httpPortMap[confHttpPort]

	basePath := server.BasePath
	if basePath == "" {
		basePath = defaultBasePath
	}

	httpHost, err := endpointHost(ctx, confHost)
	if err != nil {
		return "", err
	}
	endpointUrl := url.URL{
		Scheme: server.Protocol,
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
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
		confWebuiHost     = "webui.http.host"
	)

	httpPortMap, err := ctx.Config.Get(confWebuiHttpPort)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", confWebuiHttpPort, err)
	}
	httpPort := httpPortMap[confWebuiHttpPort]

	httpHost, err := endpointHost(ctx, confWebuiHost)
	if err != nil {
		return "", err
	}

	endpointUrl := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
		Path:   defaultBasePath,
	}

	return endpointUrl.String(), nil
}

func endpointHost(ctx *Context, hostConfigKey string) (string, error) {
	hostMap, err := ctx.Config.Get(hostConfigKey)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", hostConfigKey, err)
	}
	host := fmt.Sprint(hostMap[hostConfigKey])

	host = strings.TrimSpace(host)

	return host, nil
}
