package common

import (
	"os"
	"testing"
)

func TestSuggestListModels(t *testing.T) {
	tests := []struct {
		name                    string
		instanceName            string
		incompatibleModelsCount int
		activeEngine            string
		want                    string
	}{
		{
			name:                    "single incompatible model, not a snap",
			instanceName:            "",
			incompatibleModelsCount: 1,
			activeEngine:            "intel-gpu",
			want:                    `Hint: There is 1 other model which is not compatible with the active intel-gpu engine. Run "<snap-instance-name> models --all" to list them.`,
		},
		{
			name:                    "multiple incompatible models, not a snap",
			instanceName:            "",
			incompatibleModelsCount: 3,
			activeEngine:            "intel-gpu",
			want:                    `Hint: There are 3 other models which are not compatible with the active intel-gpu engine. Run "<snap-instance-name> models --all" to list them.`,
		},
		{
			name:                    "single incompatible model, running as a snap",
			instanceName:            "example-snap",
			incompatibleModelsCount: 1,
			activeEngine:            "cpu",
			want:                    `Hint: There is 1 other model which is not compatible with the active cpu engine. Run "example-snap models --all" to list them.`,
		},
		{
			name:                    "multiple incompatible models, running as a snap",
			instanceName:            "example-snap",
			incompatibleModelsCount: 2,
			activeEngine:            "cpu",
			want:                    `Hint: There are 2 other models which are not compatible with the active cpu engine. Run "example-snap models --all" to list them.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.instanceName == "" {
				if err := os.Unsetenv("SNAP_INSTANCE_NAME"); err != nil {
					t.Fatalf("unsetting SNAP_INSTANCE_NAME: %v", err)
				}
			} else {
				if err := os.Setenv("SNAP_INSTANCE_NAME", tt.instanceName); err != nil {
					t.Fatalf("setting SNAP_INSTANCE_NAME: %v", err)
				}
				t.Cleanup(func() { _ = os.Unsetenv("SNAP_INSTANCE_NAME") })
			}

			got := SuggestListModels(tt.incompatibleModelsCount, tt.activeEngine)
			if got != tt.want {
				t.Errorf("SuggestListModels() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}
