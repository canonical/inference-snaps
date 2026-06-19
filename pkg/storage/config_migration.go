package storage

import (
	"fmt"
	"strings"
)

func migrateConfig(c Config) error {
	if err := migratePassthroughEnv(c); err != nil {
		return fmt.Errorf("migrating passthrough environment variables: %v", err)
	}
	return nil
}

// migratePassthroughEnv migrates environment variables configurations
// See https://github.com/canonical/inference-snaps-cli/pull/361
func migratePassthroughEnv(c Config) error {
	deprecatedPrefix := "passthrough.environment"
	newPrefix := "env"

	// Get deprecated configurations
	values, err := c.Get(deprecatedPrefix)
	if err != nil {
		return err
	}

	if len(values) == 0 {
		return nil // Nothing to migrate
	}

	// Set new configurations
	for fullKey, value := range values {
		key := newPrefix + "." + strings.TrimPrefix(fullKey, deprecatedPrefix+".")
		if err := c.Set(key, fmt.Sprint(value), UserConfig); err != nil {
			return fmt.Errorf("setting %s: %s", key, err)
		}
	}

	// Unset deprecated configurations
	if err := c.Unset(deprecatedPrefix, UserConfig); err != nil {
		return fmt.Errorf("unsetting %s: %s", deprecatedPrefix, err)
	}

	return nil
}
