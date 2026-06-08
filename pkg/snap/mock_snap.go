package snap

import "fmt"

type mockSnap struct {
	installComponentFn func(name string) error
}

// Mock returns a no-op Snap suitable for tests that don't exercise InstallComponent.
func Mock() Snap {
	return &mockSnap{}
}

// MockWithInstall returns a Snap whose InstallComponent calls fn.
// Use a closure to simulate sequences of errors across multiple calls.
func MockWithInstall(fn func(name string) error) Snap {
	return &mockSnap{installComponentFn: fn}
}

func (c *mockSnap) Restart(service ...string) error {
	if len(service) == 0 {
		fmt.Println("[mock] Restarting all services")
		return nil
	}
	fmt.Println("[mock] Restarting services:", service)
	return nil
}

func (c *mockSnap) InstanceName() string {
	return "mock-snap"
}

func (c *mockSnap) InstallComponent(name string) error {
	if c.installComponentFn != nil {
		return c.installComponentFn(name)
	}
	return nil
}
