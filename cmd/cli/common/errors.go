package common

import "errors"

// Error types that can be checked higher up in the caller chain
var (
	ErrPermissionDenied = errors.New("permission denied, try again with sudo")
	ErrNoActiveEngine   = errors.New("no active engine")
)

// Strings that are commonly used in error chains, but should not be used as error types
const (
	LookingUpActiveEngine = "looking up active engine"
)
