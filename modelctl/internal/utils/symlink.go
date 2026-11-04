package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RemoveTempSymlink(path string) (success bool, err error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // nothing to remove, consider it a success
		}
		return false, fmt.Errorf("stat %s: %v", path, err)
	}

	// Check if the path is within /tmp before removing
	if withinTmp, err := pathWithinTmp(path); err != nil {
		return false, fmt.Errorf("checking if path is within /tmp: %v", err)
	} else if !withinTmp {
		return false, fmt.Errorf("layout path outside of /tmp: %s", path)
	}

	// Only remove if it's a symlink to avoid accidentally deleting other files
	if info.Mode()&os.ModeSymlink != 0 {
		if err := os.Remove(path); err != nil {
			return false, fmt.Errorf("removing file %s: %v", path, err)
		}
	} else {
		return false, fmt.Errorf("not a symlink: %s", path)
	}

	return true, nil // symlink removed successfully
}

func CreateTempSymlink(target, link string) error {
	// Reject any layout path that is outside of /tmp
	if withinTmp, err := pathWithinTmp(link); err != nil {
		return fmt.Errorf("checking if path is within /tmp: %v", err)
	} else if !withinTmp {
		return fmt.Errorf("layout path outside of /tmp: %s", link)
	}

	// Create directory tree for the link
	if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
		return fmt.Errorf("creating directory for symlink %q: %v", link, err)
	}

	// Remove existing symlink if it exists
	if _, err := RemoveTempSymlink(link); err != nil {
		return fmt.Errorf("removing existing symlink at %q: %v", link, err)
	}

	// Create new symlink
	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("creating symlink from %q to %q: %v", link, target, err)
	}

	return nil
}

func pathWithinTmp(path string) (bool, error) {
	link, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("getting absolute path of %s: %v", path, err)
	}
	if strings.HasPrefix(link, "/tmp") {
		return true, nil
	}
	return false, nil
}
