package storage

type mockCache struct {
	activeEngine string
	activeModel  string
}

func NewMockCache() Cache {
	return &mockCache{}
}

func (c *mockCache) SetActiveEngine(engine string) error {
	c.activeEngine = engine
	return nil
}

func (c *mockCache) SetActiveModel(model string) error {
	c.activeModel = model
	return nil
}

// GetActiveEngine returns the currently active engine name, or an empty string if none is set
func (c *mockCache) GetActiveEngine() (string, error) {
	return c.activeEngine, nil
}

// GetActiveModel returns the currently active model name, or an empty string if none is set
func (c *mockCache) GetActiveModel() (string, error) {
	return c.activeModel, nil
}
