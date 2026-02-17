package common

import (
	"fmt"
	"os"

	"github.com/canonical/go-snapctl/env"
)

func ComponentInstalled(component string) (bool, error) {
	// Check in /snap/$SNAP_INSTANCE_NAME/components/$SNAP_REVISION if component is mounted
	directoryPath := fmt.Sprintf("/snap/%s/components/%s/%s", env.SnapInstanceName(), env.SnapRevision(), component)

	info, err := os.Stat(directoryPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, fmt.Errorf("error checking component directory %q: %v", component, err)
		}
	} else {
		if info.IsDir() {
			return true, nil
		} else {
			return false, fmt.Errorf("component %q exists but is not a directory", component)
		}
	}
}
