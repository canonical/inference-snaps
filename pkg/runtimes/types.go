package runtimes

import (
	"slices"

	"github.com/canonical/inference-snaps-cli/v2/pkg/engines"
)

type Manifest struct {
	Name        string                    `json:"name" yaml:"name"`
	Servers     map[string]Server         `json:"servers" yaml:"servers"`
	Environment []string                  `json:"environment" yaml:"environment"`
	Layout      map[string]engines.Layout `yaml:"layout"`
	Components  []string                  `json:"components" yaml:"components"`
}

type Server struct {
	Protocol  string `json:"protocol" yaml:"protocol"`
	BasePath  string `json:"base-path" yaml:"base-path"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

func (s *Server) IsUnixProtocol() bool {
	unixProtocols := []string{
		ProtocolHttpUnix,
		ProtocolHttpsUnix,
		ProtocolWebSocketUnix,
		ProtocolWebSocketSecureUnix,
	}
	return slices.Contains(unixProtocols, s.Protocol)
}

func (s *Server) UnixSocketName() string {
	if s.Namespace != "" {
		return s.Namespace + ".sock"
	}
	return "server.sock"
}
