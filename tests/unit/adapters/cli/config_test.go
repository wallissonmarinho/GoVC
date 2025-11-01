package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
)

// TestCLIConfigFromContext tests CLIConfig creation from context
func TestCLIConfigFromContext(t *testing.T) {
	t.Run("NewCLIConfigFromContext_Valid", func(t *testing.T) {
		workers := 4
		saveLogs := true
		inputDir := "/input"

		config, err := cli.NewCLIConfigFromContext(workers, saveLogs, inputDir)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "/input", config.GetInputDir())
		assert.Equal(t, 4, config.GetWorkers())
		assert.True(t, config.SaveLogsEnabled())
	})

	t.Run("NewCLIConfigFromContext_InvalidDir", func(t *testing.T) {
		workers := 2
		saveLogs := false
		inputDir := ""

		config, err := cli.NewCLIConfigFromContext(workers, saveLogs, inputDir)

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("NewCLIConfigFromContext_DefaultWorkers", func(t *testing.T) {
		workers := 0
		saveLogs := false
		inputDir := "/input"

		config, err := cli.NewCLIConfigFromContext(workers, saveLogs, inputDir)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, 1, config.GetWorkers())
	})

	t.Run("NewCLIConfigFromContext_OutputPath", func(t *testing.T) {
		workers := 2
		saveLogs := true
		inputDir := "/videos"

		config, err := cli.NewCLIConfigFromContext(workers, saveLogs, inputDir)

		assert.NoError(t, err)
		assert.Contains(t, config.GetOutputDir(), "/videos")
		assert.Contains(t, config.GetOutputDir(), "mp4")
	})

	t.Run("NewCLIConfigFromContext_SaveLogsFlag", func(t *testing.T) {
		workers := 2
		inputDir := "/input"

		// Test with saveLogs true
		config1, _ := cli.NewCLIConfigFromContext(workers, true, inputDir)
		assert.True(t, config1.SaveLogsEnabled())

		// Test with saveLogs false
		config2, _ := cli.NewCLIConfigFromContext(workers, false, inputDir)
		assert.False(t, config2.SaveLogsEnabled())
	})
}

// TestLoggerReporter tests progress reporting
func TestLoggerReporter(t *testing.T) {
	t.Run("NewLoggerReporter", func(t *testing.T) {
		reporter := cli.NewLoggerReporter()

		assert.NotNil(t, reporter)
	})

	t.Run("LoggerReporter_ReportConversionStart", func(t *testing.T) {
		reporter := cli.NewLoggerReporter()

		// Should not panic
		assert.NotPanics(t, func() {
			reporter.ReportConversionStart("test.mkv", false)
		})
	})

	t.Run("LoggerReporter_ReportConversionFinish", func(t *testing.T) {
		reporter := cli.NewLoggerReporter()

		// Should not panic
		assert.NotPanics(t, func() {
			reporter.ReportConversionFinish("test.mkv", "/output/test.mp4", true)
		})
	})

	t.Run("LoggerReporter_ReportError", func(t *testing.T) {
		reporter := cli.NewLoggerReporter()

		// Should not panic
		assert.NotPanics(t, func() {
			reporter.ReportError("Test error message")
		})
	})

	t.Run("LoggerReporter_ReportProgress", func(t *testing.T) {
		reporter := cli.NewLoggerReporter()

		// Should not panic
		assert.NotPanics(t, func() {
			progress := map[string]float64{
				"test.mkv": 50.0,
			}
			reporter.ReportProgress(progress, true)
		})
	})
}

// TestCommandExecutor tests command execution
func TestCommandExecutor(t *testing.T) {
	t.Run("NewCommandExecutor", func(t *testing.T) {
		executor := cli.NewCommandExecutor()

		assert.NotNil(t, executor)
	})

	t.Run("CommandExecutor_Register", func(t *testing.T) {
		executor := cli.NewCommandExecutor()
		mockCmd := cli.NewConvertCommand(nil, "test")

		// Should not panic
		assert.NotPanics(t, func() {
			executor.Register("test-cmd", mockCmd)
		})
	})

	t.Run("CommandExecutor_Execute_NotRegistered", func(t *testing.T) {
		executor := cli.NewCommandExecutor()

		err := executor.Execute("non-existent")

		assert.Error(t, err)
	})
}
