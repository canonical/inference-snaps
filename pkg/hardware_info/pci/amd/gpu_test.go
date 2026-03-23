package amd

import (
	"strings"
	"testing"

	"github.com/canonical/inference-snaps-cli/pkg/types"
)

func TestVRam(t *testing.T) {
	tests := []struct {
		name          string
		device        types.PciDevice
		globalRootDir string
		expected      uint64
		shouldErr     bool
	}{
		{
			name:          "valid vram read",
			device:        types.PciDevice{Slot: "0000:03:00.0"},
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			expected:      8573157376,
			shouldErr:     false,
		},
		{
			name:          "invalid path",
			device:        types.PciDevice{Slot: "9999:99:99.9"},
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			shouldErr:     true,
		},
		{
			name:          "valid vram read",
			device:        types.PciDevice{Slot: "0000:c4:00.0"},
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			expected:      8589934592,
			shouldErr:     false,
		},
		{
			name:          "invalid path",
			device:        types.PciDevice{Slot: "9999:99:99.9"},
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			shouldErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vRam(tt.device, tt.globalRootDir)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("expected non-nil vram value")
			}
			if *got != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, *got)
			}
		})
	}
}

func TestGetAmdGpuPciSlot(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		globalRootDir string
		expected      string
		shouldErr     bool
		errContains   string
	}{
		{
			name:          "valid input with existing render minor",
			input:         "drm_render_minor 129",
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			expected:      "0000:03:00.0",
			shouldErr:     false,
		},
		{
			name:          "invalid format - missing value",
			input:         "drm_render_minor",
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			shouldErr:     true,
			errContains:   "unexpected format for drm_render_minor",
		},
		{
			name:          "invalid format - too many parts",
			input:         "drm_render_minor 128 extra",
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			shouldErr:     true,
			errContains:   "unexpected format for drm_render_minor",
		},
		{
			name:          "invalid symlink path",
			input:         "drm_render_minor 999",
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			shouldErr:     true,
		},
		{
			name:          "valid input with existing render minor",
			input:         "drm_render_minor 128",
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			expected:      "0000:c4:00.0",
			shouldErr:     false,
		},
		{
			name:          "invalid format - missing value",
			input:         "drm_render_minor",
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			shouldErr:     true,
			errContains:   "unexpected format for drm_render_minor",
		},
		{
			name:          "invalid format - too many parts",
			input:         "drm_render_minor 128 extra",
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			shouldErr:     true,
			errContains:   "unexpected format for drm_render_minor",
		},
		{
			name:          "invalid symlink path",
			input:         "drm_render_minor 999",
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			shouldErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAmdGpuPciSlot(tt.input, tt.globalRootDir)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result: %q)", got)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetGfxTargetVersion(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		errContains   string
		expectFailure bool
	}{
		{
			name:     "valid gfx target version",
			input:    "gfx_target_version 110502",
			expected: "gfx1152",
		},
		{
			name:          "invalid zero value",
			input:         "gfx_target_version 0",
			errContains:   "gfx_target_version is invalid for this device",
			expectFailure: true,
		},
		{
			name:          "unexpected format missing value",
			input:         "gfx_target_version",
			errContains:   "unexpected format for gfx_target_version",
			expectFailure: true,
		},
		{
			name:          "unexpected major format non numeric",
			input:         "gfx_target_version ab1234",
			errContains:   "error parsing major version from gfx_target_version",
			expectFailure: true,
		},
		{
			name:          "unexpected minor format non numeric",
			input:         "gfx_target_version 12ab34",
			errContains:   "error parsing minor version from gfx_target_version",
			expectFailure: true,
		},
		{
			name:          "unexpected revision format non numeric",
			input:         "gfx_target_version 1234ab",
			errContains:   "error parsing revision from gfx_target_version",
			expectFailure: true,
		},
		{
			name:          "unexpected short numeric format",
			input:         "gfx_target_version 12345",
			errContains:   "gfx_target_version has an unexpected format",
			expectFailure: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGfxTargetVersion(tt.input)
			t.Logf("input=%q got=%q err=%v", tt.input, got, err)

			if tt.expectFailure {
				if err == nil {
					t.Fatalf("expected error, got nil (result: %q)", got)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGpuPropertiesFromDir(t *testing.T) {
	tests := []struct {
		name           string
		device         types.PciDevice
		globalRootDir  []string // variadic arg
		shouldErr      bool
		checkVram      bool
		checkMicroArch bool
	}{
		{
			name:           "with specified root directory",
			device:         types.PciDevice{Slot: "0000:03:00.0"},
			globalRootDir:  []string{"../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/"},
			shouldErr:      false,
			checkVram:      true,
			checkMicroArch: true,
		},
		{
			name:          "invalid pciSlot with specified machine",
			device:        types.PciDevice{Slot: "9999:99:99.9"},
			globalRootDir: []string{"../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/"},
			shouldErr:     true,
		},
		{
			name:           "with specified root directory",
			device:         types.PciDevice{Slot: "0000:c4:00.0"},
			globalRootDir:  []string{"../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/"},
			shouldErr:      false,
			checkVram:      true,
			checkMicroArch: true,
		},
		{
			name:          "invalid pciSlot with specified machine",
			device:        types.PciDevice{Slot: "9999:99:99.9"},
			globalRootDir: []string{"../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/"},
			shouldErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var properties map[string]string
			var err error

			if len(tt.globalRootDir) > 0 {
				properties, err = gpuPropertiesFromDir(tt.device, tt.globalRootDir[0])
			} else {
				properties, err = gpuProperties(tt.device)
			}

			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkVram {
				if vram, ok := properties["vram"]; !ok || vram == "" {
					t.Fatalf("expected vram property to be set")
				}
			}

			if tt.checkMicroArch {
				if microarch, ok := properties["microarchitecture"]; !ok || microarch == "" {
					t.Fatalf("expected microarchitecture property to be set")
				}
			}
		})
	}
}

func TestGfxArchitecture(t *testing.T) {
	tests := []struct {
		name          string
		device        types.PciDevice
		globalRootDir string
		expected      string
		shouldErr     bool
		errContains   string
	}{
		{
			name:          "valid case with matching pci slot and valid gfx_target_version",
			device:        types.PciDevice{Slot: "0000:03:00.0"},
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			expected:      "gfx1032",
			shouldErr:     false,
		},
		{
			name:          "invalid nodes directory",
			device:        types.PciDevice{Slot: "0000:03:00.0"},
			globalRootDir: "/nonexistent/path/",
			shouldErr:     true,
		},
		{
			name:          "no matching node for pci slot",
			device:        types.PciDevice{Slot: "9999:99:99.9"},
			globalRootDir: "../../../../test_data/machines/hp-zbook-i712850HX+RadeonPROW6600M/machine-root/",
			shouldErr:     true,
		},
		{
			name:          "valid case with matching pci slot and valid gfx_target_version",
			device:        types.PciDevice{Slot: "0000:c4:00.0"},
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			expected:      "gfx1152",
			shouldErr:     false,
		},
		{
			name:          "no matching node for pci slot",
			device:        types.PciDevice{Slot: "9999:99:99.9"},
			globalRootDir: "../../../../test_data/machines/lenovo-thinkpad-p16s/machine-root/",
			shouldErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gfxArchitecture(tt.device, tt.globalRootDir)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result: %q)", got)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
