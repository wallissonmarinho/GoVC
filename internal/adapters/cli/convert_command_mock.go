package cli

import (
	"github.com/stretchr/testify/mock"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// MockConvertCommand is a mock implementation of ServiceCommand for testing
type MockConvertCommand struct {
	mock.Mock
}

var _ ports.ServiceCommand = (*MockConvertCommand)(nil)

// Execute runs the executor
func (m *MockConvertCommand) Execute() error {
	args := m.Called()
	return args.Error(0)
}

// Name returns the command name
func (m *MockConvertCommand) Name() string {
	args := m.Called()
	return args.String(0)
}
