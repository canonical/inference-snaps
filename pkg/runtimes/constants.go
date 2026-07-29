package runtimes

const (
	// Protocols
	ProtocolHttp                = "http"
	ProtocolHttps               = "https"
	ProtocolHttpUnix            = "http+unix"
	ProtocolHttpsUnix           = "https+unix"
	ProtocolWebSocket           = "ws"
	ProtocolWebSocketSecure     = "wss"
	ProtocolWebSocketUnix       = "ws+unix"
	ProtocolWebSocketSecureUnix = "wss+unix"
	// HTTP configuration keys
	HttpHostConfKey       = "http.host"
	HttpPortConfKey       = "http.port"
	HttpUnixSocketConfKey = "http.unix-socket"
	// WebSocket configuration keys
	WebSocketHostConfKey       = "ws.host"
	WebSocketPortConfKey       = "ws.port"
	WebSocketUnixSocketConfKey = "ws.unix-socket"

	// OpenAI server key for the runtime manifest.
	// This is used to identify OpenAI-compatible servers.
	OpenAiServerType = "openai"
)
