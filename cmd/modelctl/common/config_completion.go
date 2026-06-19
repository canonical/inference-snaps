package common

import (
	"sort"
	"strings"

	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

func CompleteConfigKeys(config storage.Config, toComplete string, appendEquals bool, excludedKeys map[string]struct{}) []string {
	values, err := config.GetAll()
	if err != nil {
		return nil
	}

	completions := make([]string, 0, len(values))
	for key := range values {
		if _, isExcluded := excludedKeys[key]; isExcluded {
			continue
		}

		candidate := key
		if appendEquals {
			candidate += "="
		}

		if strings.HasPrefix(candidate, toComplete) {
			completions = append(completions, candidate)
		}
	}

	sort.Strings(completions)
	return completions
}
