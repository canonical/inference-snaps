package nvidia

import (
	"testing"
)

func TestVRam(t *testing.T) {
	tests := []struct {
		name      string
		testInput *string
		expected  *uint64
		err       error
		shouldErr bool
	}{
		{
			name:      "converts MiB to bytes",
			testInput: new("4096 MiB"),
			expected:  new(uint64(4096 * 1024 * 1024)),
		},
		{
			name:      "returns nil for unavailable VRAM",
			testInput: new("[N/A]"),
			expected:  nil,
		},
		{
			name:      "reports parsing errors",
			testInput: new("not-a-number MiB"),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVramAmount(tt.testInput)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expected == nil {
				if got != nil {
					t.Fatalf("expected nil, got %d", *got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected non-nil result")
			}
			if *got != *tt.expected {
				t.Fatalf("expected %d, got %d", *tt.expected, *got)
			}
		})
	}
}
