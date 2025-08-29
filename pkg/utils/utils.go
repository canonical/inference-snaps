package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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

// FmtBytes converts bytes to a printable string with unit
func FmtBytes(bytes uint64) string {
	if bytes > 1024*1024*1024*1024 {
		return fmt.Sprintf("%.1fTB", float64(bytes)/1024/1024/1024/1024)
	} else if bytes > 1024*1024*1024 {
		return fmt.Sprintf("%.1fGB", float64(bytes)/1024/1024/1024)
	} else if bytes > 1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(bytes)/1024/1024)
	} else if bytes > 1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%d", bytes)
}

func StringToBytes(sizeString string) (uint64, error) {
	var sizeBytes uint64
	var scaling uint64 = 1
	var err error

	if strings.HasSuffix(sizeString, "G") {
		sizeString = strings.TrimSuffix(sizeString, "G")
		scaling = 1024 * 1024 * 1024
	} else if strings.HasSuffix(sizeString, "M") {
		sizeString = strings.TrimSuffix(sizeString, "M")
		scaling = 1024 * 1024
	}

	sizeBytes, err = strconv.ParseUint(sizeString, 10, 64)
	if err != nil {
		return 0, err
	}
	sizeBytes = sizeBytes * scaling

	return sizeBytes, nil
}

// SplitPathIntoDirectories takes a file path and returns a slice of strings containing the individual directory names that makes up the path
func SplitPathIntoDirectories(p string) []string {
	var parts []string
	for {
		dir, file := filepath.Split(p)
		if file != "" {
			parts = append([]string{file}, parts...)
		}
		if dir == "" || dir == "/" || dir == "\\" { // Handle root directory and empty path
			break
		}
		p = strings.TrimSuffix(dir, string(filepath.Separator)) // Remove trailing separator
	}
	return parts
}

// IsPrimitive takes any variable and returns true if the underlying type is a primitive type
func IsPrimitive(p interface{}) bool {
	switch reflect.TypeOf(p).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Bool, reflect.String:
		return true
	default:
		return false
	}
}

func SubDirectories(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}
	return directories, nil
}
