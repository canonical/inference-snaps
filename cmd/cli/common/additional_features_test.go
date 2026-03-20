package common

import (
	"os"
	"testing"
)

func TestAdditionalFeatures(t *testing.T) {
	const additionalFeaturesEnv = "ADDITIONAL_FEATURES"

	tests := map[string]struct {
		value            string
		setEnv           bool
		expectedFeatures additionalFeatures
	}{
		"chat only": {
			value:  "chat",
			setEnv: true,
			expectedFeatures: additionalFeatures{
				Chat: true,
			},
		},
		"chat and others": {
			value:  "chat,ui",
			setEnv: true,
			expectedFeatures: additionalFeatures{
				Chat: true,
			},
		},
		"chat with whitespaces": {
			value:  " chat ,ui",
			setEnv: true,
			expectedFeatures: additionalFeatures{
				Chat: true,
			},
		},
		"no chat": {
			value:  "ui",
			setEnv: true,
			expectedFeatures: additionalFeatures{
				Chat: false,
			},
		},
		"missing value": {
			setEnv:           false,
			expectedFeatures: additionalFeatures{},
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			if testData.setEnv {
				t.Setenv(additionalFeaturesEnv, testData.value)
			} else {
				originalValue, hadOriginalValue := os.LookupEnv(additionalFeaturesEnv)

				if err := os.Unsetenv(additionalFeaturesEnv); err != nil {
					t.Fatalf("error unsetting %s: %v", additionalFeaturesEnv, err)
				}

				t.Cleanup(func() {
					if hadOriginalValue {
						_ = os.Setenv(additionalFeaturesEnv, originalValue)
					}
				})
			}

			features := AdditionalFeatures()
			if features != testData.expectedFeatures {
				t.Errorf("returned %+v, expected %+v", features, testData.expectedFeatures)
			}
		})
	}
}
