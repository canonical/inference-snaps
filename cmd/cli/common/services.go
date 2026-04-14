package common

import (
	"fmt"
	"strings"

	"github.com/canonical/go-snapctl"
)

func ServiceStatuses() (map[string]string, error) {
	services, err := snapctl.Services().Run()
	if err != nil {
		return nil, fmt.Errorf("getting list of services: %v", err)
	}
	statuses := make(map[string]string)
	for name, service := range services {
		// The service name is in the format <snap-name>.<service-app>, we only want the service-app part.
		_, serviceApp, found := strings.Cut(name, ".")
		if !found {
			return nil, fmt.Errorf("unexpected service name format: %q", name)
		}
		// Append the service status exactly as snapd reports it. Often this is in the host system language, see bug:
		// https://bugs.launchpad.net/snapd/+bug/2137543
		statuses[serviceApp] = service.Current
	}
	return statuses, nil
}
