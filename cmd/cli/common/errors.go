package common

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied, try again with sudo")
)
