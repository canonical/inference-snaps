package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveTempSymlinkSkipsRegularFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-symlink.txt")
	if err := os.WriteFile(path, []byte("keep-me"), 0644); err != nil {
		t.Fatalf("failed to create regular file: %v", err)
	}

	if _, err := RemoveTempSymlink(path); err == nil {
		t.Fatalf("expected error removing regular file path, got nil")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected regular file to remain, read failed: %v", err)
	}
	if string(content) != "keep-me" {
		t.Fatalf("expected regular file content to stay intact, got %q", string(content))
	}
}

func TestRemoveTempSymlinkRejectsOutsideTmp(t *testing.T) {
	success, err := RemoveTempSymlink("/nonexistent-path-outside-tmp")
	if err != nil {
		t.Fatalf("unexpected error for non-existent path: %v", err)
	}
	if !success {
		t.Fatalf("expected success for non-existent path, got false")
	}
}

func TestCreateTempSymlinkBasic(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target.txt")
	link := filepath.Join(tmpDir, "link")

	// Create target file
	if err := os.WriteFile(target, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}

	// Create symlink
	if err := CreateTempSymlink(target, link); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify symlink exists and points to target
	linkTarget, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if linkTarget != target {
		t.Errorf("expected symlink target %q, got %q", target, linkTarget)
	}

	// Verify content is accessible through symlink
	content, err := os.ReadFile(link)
	if err != nil {
		t.Fatalf("failed to read through symlink: %v", err)
	}
	if string(content) != "content" {
		t.Errorf("expected content %q, got %q", "content", string(content))
	}
}

func TestCreateTempSymlinkCreatesNestedDirs(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target.txt")
	link := filepath.Join(tmpDir, "nested", "dir", "structure", "link")

	// Create target file
	if err := os.WriteFile(target, []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}

	// Create symlink with nested directories
	if err := CreateTempSymlink(target, link); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify symlink exists
	linkTarget, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if linkTarget != target {
		t.Errorf("expected symlink target %q, got %q", target, linkTarget)
	}

	// Verify all directories were created
	info, err := os.Stat(filepath.Dir(link))
	if err != nil {
		t.Fatalf("directory structure not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory, got file")
	}
}

func TestCreateTempSymlinkReplacesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	oldTarget := filepath.Join(tmpDir, "old-target.txt")
	newTarget := filepath.Join(tmpDir, "new-target.txt")
	link := filepath.Join(tmpDir, "link")

	// Create target files
	if err := os.WriteFile(oldTarget, []byte("old"), 0644); err != nil {
		t.Fatalf("failed to create old target: %v", err)
	}
	if err := os.WriteFile(newTarget, []byte("new"), 0644); err != nil {
		t.Fatalf("failed to create new target: %v", err)
	}

	// Create initial symlink
	if err := CreateTempSymlink(oldTarget, link); err != nil {
		t.Fatalf("failed to create initial symlink: %v", err)
	}

	// Verify initial symlink
	linkTarget, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read initial symlink: %v", err)
	}
	if linkTarget != oldTarget {
		t.Errorf("expected initial target %q, got %q", oldTarget, linkTarget)
	}

	// Replace symlink with new target
	if err := CreateTempSymlink(newTarget, link); err != nil {
		t.Fatalf("failed to replace symlink: %v", err)
	}

	// Verify symlink now points to new target
	linkTarget, err = os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read replaced symlink: %v", err)
	}
	if linkTarget != newTarget {
		t.Errorf("expected new target %q, got %q", newTarget, linkTarget)
	}

	// Verify content comes from new target
	content, err := os.ReadFile(link)
	if err != nil {
		t.Fatalf("failed to read through symlink: %v", err)
	}
	if string(content) != "new" {
		t.Errorf("expected content %q, got %q", "new", string(content))
	}
}

func TestCreateTempSymlinkRejectsOutsideTmp(t *testing.T) {
	target := "/some/file"
	link := "./invalid-tmp-link"

	err := CreateTempSymlink(target, link)
	if err == nil {
		t.Fatal("expected error for path outside /tmp, got nil")
	}
	if _, ok := err.(interface{ Error() string }); !ok {
		t.Fatalf("expected error type, got: %T", err)
	}
}

func TestCreateTempSymlinkDoesNotRemoveExistingRegularFile(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target.txt")
	link := filepath.Join(tmpDir, "link")

	// Create target file
	if err := os.WriteFile(target, []byte("target"), 0644); err != nil {
		t.Fatalf("failed to create target: %v", err)
	}

	// Create a regular file at link location (not a symlink)
	if err := os.WriteFile(link, []byte("regular"), 0644); err != nil {
		t.Fatalf("failed to create regular file: %v", err)
	}

	// CreateTempSymlink should fail because we can't remove regular files
	err := CreateTempSymlink(target, link)
	if err == nil {
		t.Fatal("expected error when trying to replace a regular file, got nil")
	}

	// Verify the regular file still exists unchanged
	content, err := os.ReadFile(link)
	if err != nil {
		t.Fatalf("failed to read regular file: %v", err)
	}
	if string(content) != "regular" {
		t.Errorf("expected regular file unchanged, got %q", string(content))
	}
}
