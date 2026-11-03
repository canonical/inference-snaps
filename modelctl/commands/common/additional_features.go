package common

import (
	"os"
	"strings"
)

type additionalFeatures struct {
	Chat  bool
	WebUi bool
}

func AdditionalFeatures() additionalFeatures {
	const (
		additionalFeaturesEnv = "ADDITIONAL_FEATURES"
		featureChat           = "chat"
		featureWebUi          = "webui"
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
		case featureWebUi:
			features.WebUi = true
		}
	}

	return features
}

func ChatEnabled() bool {
	features := AdditionalFeatures()
	return features.Chat
}

func WebUiEnabled() bool {
	features := AdditionalFeatures()
	return features.WebUi
}
