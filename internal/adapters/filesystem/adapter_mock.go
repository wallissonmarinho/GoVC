package filesystem

import (
	"github.com/wallissonmarinho/GoVC/internal/core/domain"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
	"github.com/stretchr/testify/mock"
)

// MockFilesystemAdapter is a mock implementation of VideoDiscoveryPort and FileSystemPort for testing
type MockFilesystemAdapter struct {
	mock.Mock
}

// VideoDiscoveryPort methods
func (m *MockFilesystemAdapter) FindVideos(dir string) ([]*domain.Video, error) {
	args := m.Called(dir)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Video), args.Error(1)
}

func (m *MockFilesystemAdapter) CreateOutputDir(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

// FileSystemPort methods
func (m *MockFilesystemAdapter) FileExists(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockFilesystemAdapter) IsValidOutput(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockFilesystemAdapter) RemoveFile(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockFilesystemAdapter) WriteLog(logPath string, lines []string) error {
	args := m.Called(logPath, lines)
	return args.Error(0)
}

// Ensure MockFilesystemAdapter implements VideoDiscoveryPort and FileSystemPort
var _ ports.VideoDiscoveryPort = (*MockFilesystemAdapter)(nil)
var _ ports.FileSystemPort = (*MockFilesystemAdapter)(nil)
