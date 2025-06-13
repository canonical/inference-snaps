package types

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// HexInt Custom type to handle hex values
type HexInt int

// UnmarshalYAML parses a hex string into an int
func (hi *HexInt) UnmarshalYAML(value *yaml.Node) error {
	// Ignore empty string
	if value.Value == "" {
		return nil
	}

	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("expected a scalar node, got %v", value.LongTag())
	}

	// Strip 0x prefix if it exists
	hexString := strings.TrimPrefix(value.Value, "0x")

	// Parse the hex string to int
	parsed, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return fmt.Errorf("failed to parse hex value %s: %w", value.Value, err)
	}

	*hi = HexInt(parsed)
	return nil
}

func (hi *HexInt) UnmarshalJSON(data []byte) error {
	// Remove quotes
	hexString := strings.Trim(string(data), "\"")

	// Ignore empty string
	if hexString == "" {
		return nil
	}

	// Remove "0x" prefix if present
	hexString = strings.TrimPrefix(hexString, "0x")

	// Parse as base 16 integer
	val, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return fmt.Errorf("failed to parse hex value %s: %w", hexString, err)
	}
	*hi = HexInt(val)
	return nil
}

func (hi HexInt) MarshalJSON() ([]byte, error) {
	hexString := fmt.Sprintf("\"%x\"", int(hi))
	return []byte(hexString), nil
}
