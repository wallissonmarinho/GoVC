package commands

import urfavecli "github.com/urfave/cli/v2"

// MockCommandFactory is a mock for testing CommandFactory
type MockCommandFactory struct {
	BuildCommandsCalled bool
	BuildCommandsError  error
}

// BuildCommands mocks the BuildCommands method
func (m *MockCommandFactory) BuildCommands() []*urfavecli.Command {
	m.BuildCommandsCalled = true
	if m.BuildCommandsError != nil {
		return nil
	}
	return []*urfavecli.Command{
		{
			Name: "mock-command",
		},
	}
}
