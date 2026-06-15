package utils

import (
	"os"
	"testing"
)

func TestStringToBytesGigabytes(t *testing.T) {
	sizeBytes, err := StringToBytes("4G")
	if err != nil {
		t.Fatal(err)
	}
	if sizeBytes != 4*1024*1024*1024 {
		t.Fatal("incorrectly parsed size")
	}
}

func TestStringToBytesMegabytes(t *testing.T) {
	sizeBytes, err := StringToBytes("256M")
	if err != nil {
		t.Fatal(err)
	}
	if sizeBytes != 256*1024*1024 {
		t.Fatal("incorrectly parsed size")
	}
}

func TestStringToBytesTerabytes(t *testing.T) {
	_, err := StringToBytes("2T")
	if err == nil {
		t.Fatal("Terabytes should not be supported")
	}
}

func TestStringToBytesKilobytes(t *testing.T) {
	_, err := StringToBytes("1024K")
	if err == nil {
		t.Fatal("Kilobytes should not be supported")
	}
}

func TestStringToBytesUnknown(t *testing.T) {
	_, err := StringToBytes("1024A")
	if err == nil {
		t.Fatal("Unknown unit should not be parsed")
	}
}

func TestStringToBytesExponent(t *testing.T) {
	// GO only supports exponents for floats
	_, err := StringToBytes("10E4")
	if err == nil {
		t.Fatal("Exponents should not be supported")
	}
}

func TestStringToBytes(t *testing.T) {
	sizeBytes, err := StringToBytes("256")
	if err != nil {
		t.Fatal(err)
	}
	if sizeBytes != 256 {
		t.Fatal("incorrectly parsed size")
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
