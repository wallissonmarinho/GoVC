package cli

import (
	"github.com/stretchr/testify/mock"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// MockCommandExecutor is a mock implementation of CommandExecutor for testing
type MockCommandExecutor struct {
	mock.Mock
}

var _ ports.CommandExecutorPort = (*MockCommandExecutor)(nil)

// Register adds a command to the executor
func (m *MockCommandExecutor) Register(name string, cmd ports.ServiceCommand) {
	m.Called(name, cmd)
}

// Execute runs a command by name
func (m *MockCommandExecutor) Execute(name string) error {
	args := m.Called(name)
	return args.Error(0)
}
