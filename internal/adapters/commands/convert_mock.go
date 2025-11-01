package commands

import urfavecli "github.com/urfave/cli/v2"

// MockConvertCommandHandler is a mock for testing ConvertCommandHandler
type MockConvertCommandHandler struct {
	ExecuteCalled bool
	ExecuteError  error
}

// Execute mocks the Execute method
func (m *MockConvertCommandHandler) Execute(c *urfavecli.Context) error {
	m.ExecuteCalled = true
	return m.ExecuteError
}

// BuildCommand mocks the BuildCommand method
func (m *MockConvertCommandHandler) BuildCommand() *urfavecli.Command {
	return &urfavecli.Command{
		Name: "convert-mock",
	}
}
