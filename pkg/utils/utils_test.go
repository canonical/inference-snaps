package utils

import (
	"os"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{name: "gibibytes", input: "4G", want: 4 * 1024 * 1024 * 1024},
		{name: "mebibytes", input: "256M", want: 256 * 1024 * 1024},
		{name: "float mebibytes", input: "3.14M", want: 3292528}, // 3.14 * 1024*1024 = 3292528.64
		{name: "bytes", input: "256", want: 256},
		{name: "NaN", input: "NaNM", wantErr: true},
		{name: "positive infinity", input: "+infM", wantErr: true},
		{name: "negative infinity", input: "-infM", wantErr: true},
		{name: "tebibytes unsupported", input: "2T", wantErr: true},
		{name: "kibibytes unsupported", input: "1024K", wantErr: true},
		{name: "GB suffix unsupported", input: "6GB", wantErr: true},
		{name: "GiB suffix unsupported", input: "6GiB", wantErr: true},
		{name: "Gi suffix unsupported", input: "6Gi", wantErr: true},
		{name: "unknown unit", input: "1024A", wantErr: true},
		{name: "exponent notation unsupported", input: "10E4", wantErr: true},
		{name: "negative values", input: "-1024M", wantErr: true},
		{name: "overflows uint64", input: "20000000000G", wantErr: true},
		{name: "exactly 2^64 boundary", input: "18446744073709551616", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToBytes(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("StringToBytes(%q) = %d, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("StringToBytes(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("StringToBytes(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsPrimitive(t *testing.T) {
	if !IsPrimitive(1) {
		t.Fatal("int should be primitive")
	}
	if !IsPrimitive("test") {
		t.Fatal("string should be primitive")
	}
	if !IsPrimitive(true) {
		t.Fatal("boolean should be primitive")
	}
	if IsPrimitive([]string{"test"}) {
		t.Fatal("string slice should not be primitive")
	}
}

// This is for manual testing
func TestIsRootUser(t *testing.T) {
	if IsRootUser() {
		t.Log("User is root!")
	} else {
		t.Log("User is not root")
	}
}

func TestFmtBytesShortUnit(t *testing.T) {
	tests := []struct {
		value    float64
		unit     string
		expected string
	}{
		{1.0, "G", "1G"},
		{2.0, "M", "2M"},
		{1.5, "G", "1.5G"},
		{0.0, "K", "0K"},
		{1024.0, "T", "1024T"},
		{3.7, "K", "3.7K"},
	}
	for _, tt := range tests {
		got := fmtBytesShortUnit(tt.value, tt.unit)
		if got != tt.expected {
			t.Errorf("fmtBytesShortUnit(%v, %q) = %q, want %q", tt.value, tt.unit, got, tt.expected)
		}
	}
}

func TestFmtBytesShort(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{0, "0"},
		{512, "512"},
		{1024, "1K"},
		{2048, "2K"},
		{1536, "1.5K"},
		{1024 * 1024, "1M"},
		{2 * 1024 * 1024, "2M"},
		{uint64(1.5 * 1024 * 1024), "1.5M"},
		{1024 * 1024 * 1024, "1G"},
		{2 * 1024 * 1024 * 1024, "2G"},
		{uint64(1.5 * 1024 * 1024 * 1024), "1.5G"},
		{1024 * 1024 * 1024 * 1024, "1T"},
		{2 * 1024 * 1024 * 1024 * 1024, "2T"},
		{uint64(1.5 * 1024 * 1024 * 1024 * 1024), "1.5T"},
	}
	for _, tt := range tests {
		got := FmtBytesShort(tt.bytes)
		if got != tt.expected {
			t.Errorf("FmtBytesShort(%d) = %q, want %q", tt.bytes, got, tt.expected)
		}
	}
}

func TestSetEnvironmentVariables(t *testing.T) {
	defer os.Unsetenv("TEST_VAR")

	envVars := map[string]any{
		"TEST_VAR": "test-value",
	}

	err := SetEnvironmentVariables(envVars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if value := os.Getenv("TEST_VAR"); value != "test-value" {
		t.Fatalf("expected TEST_VAR to be 'test-value', got '%s'", value)
	}
}
