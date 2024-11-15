package memory

import (
	"golang.org/x/sys/unix"
)

func sysInfo() (*unix.Sysinfo_t, error) {
	var info unix.Sysinfo_t
	err := unix.Sysinfo(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
