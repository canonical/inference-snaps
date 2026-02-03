package common

import (
	"fmt"

	"github.com/canonical/go-snapctl/env"
)

func SuggestServerStartup() string {
	return "Try again when the server is ready."
}

func SuggestServerLogs() string {

	instanceName := env.SnapInstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	// TODO: get app name dynamically
	serviceName := instanceName + ".server"

	return fmt.Sprintf("Run \"snap logs %s\" to see the server logs.", serviceName)
}

func SuggestStartServer() string {

	instanceName := env.SnapInstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	// TODO: get app name dynamically
	serviceName := instanceName + ".server"

	return fmt.Sprintf("Run \"sudo snap start %s\" to start the server.", serviceName)
}

func SuggestServiceManagement() string {

	instanceName := env.SnapInstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	return fmt.Sprintf("Use \"snap logs|start|stop|restart %v\" for service management.", instanceName)
}

func SuggestEngineInfo() string {
	instanceName := env.SnapInstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	return fmt.Sprintf("Use \"%v show-engine <engine>\" for more information about an engine.", instanceName)
}
