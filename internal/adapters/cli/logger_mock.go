package cli

import (
	"github.com/stretchr/testify/mock"
)

// MockLoggerReporter is a mock implementation of ProgressReporterPort for testing
type MockLoggerReporter struct {
	mock.Mock
}

func (m *MockLoggerReporter) ReportProgress(progress map[string]float64, isComplete bool) {
	m.Called(progress, isComplete)
}

func (m *MockLoggerReporter) ReportConversionStart(filename string, hasSubs bool) {
	m.Called(filename, hasSubs)
}

func (m *MockLoggerReporter) ReportConversionFinish(filename string, outputPath string, success bool) {
	m.Called(filename, outputPath, success)
}

func (m *MockLoggerReporter) ReportError(message string) {
	m.Called(message)
}
