package common

import (
	"fmt"

	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
)

func SuggestServerStartup() string {
	return "Try again when the server is ready."
}

func SuggestServerLogs() string {

	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	// TODO: get app name dynamically
	serviceName := instanceName + ".server"

	return fmt.Sprintf("Run \"snap logs %s\" to see the server logs.", serviceName)
}

func SuggestStartServer() string {

	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	// TODO: get app name dynamically
	serviceName := instanceName + ".server"

	return fmt.Sprintf("Run \"sudo snap start %s\" to start the server.", serviceName)
}

func SuggestStartService(service string) string {

	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	serviceName := instanceName + "." + service

	return fmt.Sprintf("Run \"sudo snap start %s\" to start it.", serviceName)
}

func SuggestServiceManagement() string {

	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	return fmt.Sprintf("Use \"snap logs|start|stop|restart %v\" for service management.", instanceName)
}

func SuggestEngineInfo() string {
	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	return fmt.Sprintf("Use \"%v engine <engine>\" for more information about an engine.", instanceName)
}

func SuggestKeyNotFound(key string) string {
	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}

	return fmt.Sprintf("Use \"%s get\" to view available keys", instanceName)
}

func SuggestListModels(incompatibleModelsCount int, activeEngine string) string {
	instanceName := snap.InstanceName()
	if instanceName == "" { // not a snap
		instanceName = "<snap-instance-name>"
	}
	var hintTemplate string
	if incompatibleModelsCount == 1 {
		hintTemplate = "There is %d other model which is not compatible with the active %s engine"
	} else {
		hintTemplate = "There are %d other models which are not compatible with the active %s engine"
	}
	command := fmt.Sprintf("%s models --all", instanceName)
	return fmt.Sprintf("Hint: "+hintTemplate+". Run %q to list them.",
		incompatibleModelsCount, activeEngine, command)
}
