package storage

type mockCache struct {
	activeEngine string
}

func NewMockCache() Cache {
	return &mockCache{}
}

func (c *mockCache) SetActiveEngine(engine string) error {
	c.activeEngine = engine
	return nil
}

// GetActiveEngine returns the currently active engine name, or an empty string if none is set
func (c *mockCache) GetActiveEngine() (string, error) {
	return c.activeEngine, nil
}
