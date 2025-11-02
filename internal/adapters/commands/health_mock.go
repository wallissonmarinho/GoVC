package commands

import urfavecli "github.com/urfave/cli/v2"

// MockHealthCommandHandler is a mock for testing HealthCommandHandler
type MockHealthCommandHandler struct {
	ExecuteCalled bool
	ExecuteError  error
}

// Execute mocks the Execute method
func (m *MockHealthCommandHandler) Execute(c *urfavecli.Context) error {
	m.ExecuteCalled = true
	return m.ExecuteError
}

// BuildCommand mocks the BuildCommand method
func (m *MockHealthCommandHandler) BuildCommand() *urfavecli.Command {
	return &urfavecli.Command{
		Name: "health-mock",
	}
}
