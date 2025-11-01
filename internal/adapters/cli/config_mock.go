package cli

import (
	"github.com/stretchr/testify/mock"

	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// MockCLIConfig is a mock implementation of ConfigPort for testing
type MockCLIConfig struct {
	mock.Mock
}

func (m *MockCLIConfig) GetInputDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCLIConfig) GetOutputDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCLIConfig) GetWorkers() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockCLIConfig) SaveLogsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

// Ensure MockCLIConfig implements ConfigPort
var _ ports.ConfigPort = (*MockCLIConfig)(nil)
