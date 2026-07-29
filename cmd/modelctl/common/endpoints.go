package common

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
	"github.com/canonical/inference-snaps-cli/v2/pkg/runtimes"
)

type Entrypoints map[string]Entrypoint

type Entrypoint struct {
	Url           string `json:"url,omitempty" yaml:"url,omitempty"`
	UnixSocket    string `json:"unix-socket,omitempty" yaml:"unix-socket,omitempty"`
	UnixSocketUrl string `json:"unix-socket-url,omitempty" yaml:"-"`
}

func (e Entrypoint) MarshalYAML() (any, error) {
	type entrypointYAML Entrypoint

	// Add UnixSocketUrl to UnixSocket as annotation
	unixSocket := e.UnixSocket
	if unixSocket != "" && e.UnixSocketUrl != "" {
		unixSocket = fmt.Sprintf("%s (%s)", unixSocket, e.UnixSocketUrl)
	}

	return entrypointYAML{
		Url:        e.Url,
		UnixSocket: unixSocket,
	}, nil
}

func ServerEntrypoints(ctx *Context) (Entrypoints, error) {
	activeEngineName, err := ctx.Cache.GetActiveEngine()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", LookingUpActiveEngine, err)
	}
	if activeEngineName == "" {
		return nil, ErrNoActiveEngine
	}
	activeEngineManifest, err := engines.LoadManifest(ctx.EnginesDir, activeEngineName)
	if err != nil {
		return nil, fmt.Errorf("loading active engine manifest: %v", err)
	}

	// If the engine does not list a runtime, return no entrypoints
	if activeEngineManifest.Runtime == "" {
		return nil, nil
	}

	runtimeManifest, err := runtimes.LoadManifest(ctx.RuntimesDir, activeEngineManifest.Runtime)
	if err != nil {
		return nil, fmt.Errorf("loading runtime manifest: %v", err)
	}

	entrypoints := make(Entrypoints)

	for serverName, serverSettings := range runtimeManifest.Servers {
		var entrypoint *Entrypoint
		var err error

		switch serverSettings.Protocol {
		case runtimes.ProtocolHttp, runtimes.ProtocolHttps:
			entrypoint, err = serverHttpEntrypoint(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("constructing HTTP entrypoint: %v", err)
			}
		case runtimes.ProtocolHttpUnix, runtimes.ProtocolHttpsUnix:
			entrypoint, err = serverHttpOverUnixSocketEntrypoint(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("constructing HTTP Unix entrypoint: %v", err)
			}
		case runtimes.ProtocolWebSocket, runtimes.ProtocolWebSocketSecure:
			entrypoint, err = serverWsEntrypoint(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("constructing WebSocket entrypoint: %v", err)
			}
		case runtimes.ProtocolWebSocketUnix, runtimes.ProtocolWebSocketSecureUnix:
			entrypoint, err = serverWsOverUnixSocketEntrypoint(ctx, serverSettings)
			if err != nil {
				return nil, fmt.Errorf("constructing WebSocket Unix entrypoint: %v", err)
			}
		default:
			return nil, fmt.Errorf("unsupported protocol %q for server %q in runtime %q",
				serverSettings.Protocol, serverName, activeEngineManifest.Runtime)
		}
		entrypoints[serverName] = *entrypoint
	}

	// If builtin webui is enabled, list it as well
	if WebUiEnabled() {
		webUiUrl, err := UiServerHttpUrl(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting web UI url: %v", err)
		}
		entrypoints["webui"] = Entrypoint{Url: webUiUrl}
	}

	return entrypoints, nil
}

// fullConfKey returns the full configuration key for a given key and namespace
func fullConfKey(key string, namespace string) string {
	if namespace != "" {
		return fmt.Sprintf("%s.%s", namespace, key)
	}
	return key
}

func serverHttpEntrypoint(ctx *Context, server runtimes.Server) (*Entrypoint, error) {
	httpPort, err := getConfigString(ctx, fullConfKey(runtimes.HttpPortConfKey, server.Namespace))
	if err != nil {
		return nil, err
	}

	httpHost, err := getConfigString(ctx, fullConfKey(runtimes.HttpHostConfKey, server.Namespace))
	if err != nil {
		return nil, err
	}

	entrypointUrl := url.URL{
		Scheme: server.Protocol,
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
		Path:   server.BasePath,
	}

	return &Entrypoint{Url: entrypointUrl.String()}, nil
}

func serverHttpOverUnixSocketEntrypoint(ctx *Context, server runtimes.Server) (*Entrypoint, error) {

	unixSocket, err := getConfigString(ctx, fullConfKey(runtimes.HttpUnixSocketConfKey, server.Namespace))
	if err != nil {
		return nil, err
	}

	protocol := strings.TrimSuffix(server.Protocol, "+unix")
	unixSocketUrl := fmt.Sprintf("%s://unix%s", protocol, server.BasePath) // remove +unix suffix for URL scheme

	return &Entrypoint{
		UnixSocket:    unixSocket,
		UnixSocketUrl: unixSocketUrl,
	}, nil
}

func serverWsEntrypoint(ctx *Context, server runtimes.Server) (*Entrypoint, error) {
	wsPort, err := getConfigString(ctx, fullConfKey(runtimes.WebSocketPortConfKey, server.Namespace))
	if err != nil {
		return nil, err
	}

	wsHost, err := getConfigString(ctx, fullConfKey(runtimes.WebSocketHostConfKey, server.Namespace))
	if err != nil {
		return nil, err
	}

	entrypointUrl := url.URL{
		Scheme: server.Protocol,
		Host:   net.JoinHostPort(wsHost, fmt.Sprint(wsPort)),
		Path:   server.BasePath,
	}

	return &Entrypoint{Url: entrypointUrl.String()}, nil
}

func serverWsOverUnixSocketEntrypoint(ctx *Context, server runtimes.Server) (*Entrypoint, error) {
	unixSocket, err := getConfigString(ctx, fullConfKey(runtimes.WebSocketUnixSocketConfKey, server.Namespace))
	if err != nil {
		return nil, err
	}

	protocol := strings.TrimSuffix(server.Protocol, "+unix") // remove +unix suffix for URL scheme
	unixSocketUrl := fmt.Sprintf("%s://unix%s", protocol, server.BasePath)

	return &Entrypoint{
		UnixSocket:    unixSocket,
		UnixSocketUrl: unixSocketUrl,
	}, nil
}

func OpenAiBaseUrl(ctx *Context) (string, error) {
	entrypoints, err := ServerEntrypoints(ctx)
	if err != nil {
		return "", fmt.Errorf("getting server entrypoints: %v", err)
	}
	entrypoint, found := entrypoints[runtimes.OpenAiServerType]
	if !found {
		return "", fmt.Errorf("%q not found in server entrypoints", runtimes.OpenAiServerType)
	}
	if entrypoint.Url == "" {
		return "", fmt.Errorf("%q entrypoint does not have a URL", runtimes.OpenAiServerType)
	}
	return entrypoint.Url, nil
}

func UiServerHttpUrl(ctx *Context) (string, error) {
	const (
		confWebuiHttpPort = "webui.http.port"
		confWebuiHost     = "webui.http.host"
	)

	httpPort, err := getConfigString(ctx, confWebuiHttpPort)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", confWebuiHttpPort, err)
	}

	httpHost, err := getConfigString(ctx, confWebuiHost)
	if err != nil {
		return "", err
	}

	entrypointUrl := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(httpHost, fmt.Sprint(httpPort)),
	}

	return entrypointUrl.String(), nil
}

func getConfigString(ctx *Context, key string) (string, error) {
	valueMap, err := ctx.Config.Get(key)
	if err != nil {
		return "", fmt.Errorf("getting config %q: %v", key, err)
	}
	value := fmt.Sprint(valueMap[key])
	value = strings.TrimSpace(value)
	if value == "" || value == "<nil>" {
		return "", fmt.Errorf("config %q is not set", key)
	}
	return value, nil
}
