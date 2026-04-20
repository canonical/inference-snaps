package snap

import "fmt"

type mockSnap struct{}

func Mock() Snap {
	return &mockSnap{}
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
