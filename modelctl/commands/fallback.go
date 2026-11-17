package commands

import (
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/fallbackserver"
)

// listenAddress parses a URL and returns the host:port address.
func listenAddress(baseURL string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parsing base URL: %w", err)
	}
	if parsed.Hostname() == "" {
		return "", fmt.Errorf("base URL %q does not contain a host", baseURL)
	}

	port := parsed.Port()
	if port == "" {
		port = "80" // default for http url without a port
	}
	return net.JoinHostPort(parsed.Hostname(), port), nil
}

func startFallbackServer(address string, diagnosisErr error) error {
	message := "Inference server failed to start. Please check the logs for more information."
	switch {
	case errors.Is(diagnosisErr, common.ErrNoActiveEngine):
		message = "No active engine is set. Please set an active engine and try again."
	case errors.Is(diagnosisErr, common.ErrNoActiveModel):
		message = "No active model is set. Please set an active model and try again."
	}

	fmt.Println("Starting fallback server...")
	if err := fallbackserver.Run(address, message); err != nil {
		return fmt.Errorf("running fallback server: %w", err)
	}
	return nil
}
