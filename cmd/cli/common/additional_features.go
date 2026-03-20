package common

import (
	"os"
	"strings"
)

type additionalFeatures struct {
	Chat bool
}

func AdditionalFeatures() additionalFeatures {
	const (
		additionalFeaturesEnv = "ADDITIONAL_FEATURES"
		featureChat           = "chat"
	)

	featuresCsv, found := os.LookupEnv(additionalFeaturesEnv)
	if !found {
		return additionalFeatures{}
	}

	var features additionalFeatures
	for feature := range strings.SplitSeq(featuresCsv, ",") {
		switch strings.TrimSpace(feature) {
		case featureChat:
			features.Chat = true
		}
	}

	return features
}
