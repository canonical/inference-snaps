package utils

import (
	"encoding/json"
	"fmt"
)

// FmtPretty converts any interface to JSON with indentation, for use in logging where better readability is required. Errors are ignored.
func FmtPretty(v interface{}) string {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// Ignore error
	}
	return string(jsonData)
}

// FmtGigabytes converts bytes to a printable string of gigabytes, rounded to the closest integer.
func FmtGigabytes(bytes uint64) string {
	return fmt.Sprintf("%.0fGB", float64(bytes)/1024/1024/1024)
}
