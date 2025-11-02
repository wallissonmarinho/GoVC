package cli

import (
	"github.com/stretchr/testify/mock"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// MockHealthCommand is a mock implementation of ServiceCommand for testing
type MockHealthCommand struct {
	mock.Mock
}

var _ ports.ServiceCommand = (*MockHealthCommand)(nil)

// Execute runs the executor
func (m *MockHealthCommand) Execute() error {
	args := m.Called()
	return args.Error(0)
}

// Name returns the command name
func (m *MockHealthCommand) Name() string {
	args := m.Called()
	return args.String(0)
}
