package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/commands"
)

// TestConvertCommandHandler tests the convert command handler
func TestConvertCommandHandler(t *testing.T) {
	t.Run("NewConvertCommandHandler", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		assert.NotNil(t, handler)
	})

	t.Run("BuildCommand_ReturnsValidCommand", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		assert.NotNil(t, cmd)
		assert.Equal(t, "convert", cmd.Name)
		assert.Equal(t, "Convert MKV files to MP4", cmd.Usage)
		assert.NotEmpty(t, cmd.Flags)
		assert.NotNil(t, cmd.Action)
	})

	t.Run("BuildCommand_HasWorkersFlagAlias", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		// Find the workers flag
		var workersFlag *urfavecli.IntFlag
		for _, flag := range cmd.Flags {
			if intFlag, ok := flag.(*urfavecli.IntFlag); ok && intFlag.Name == "workers" {
				workersFlag = intFlag
				break
			}
		}

		assert.NotNil(t, workersFlag)
		assert.Contains(t, workersFlag.Aliases, "p")
	})

	t.Run("BuildCommand_HasLogsFlag", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		// Find the logs flag
		var logsFlag *urfavecli.BoolFlag
		for _, flag := range cmd.Flags {
			if boolFlag, ok := flag.(*urfavecli.BoolFlag); ok && boolFlag.Name == "logs" {
				logsFlag = boolFlag
				break
			}
		}

		assert.NotNil(t, logsFlag)
		assert.False(t, logsFlag.Value) // Default should be false (delete logs)
	})

	t.Run("BuildCommand_HasDescription", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		assert.NotEmpty(t, cmd.Description)
		assert.Contains(t, cmd.Description, "MKV")
		assert.Contains(t, cmd.Description, "MP4")
	})

	t.Run("BuildCommand_FlagsHaveCorrectTypes", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		flagCount := len(cmd.Flags)
		assert.Greater(t, flagCount, 0, "Command should have flags")

		// Check for specific flag types
		hasIntFlag := false
		hasBoolFlag := false

		for _, flag := range cmd.Flags {
			switch flag.(type) {
			case *urfavecli.IntFlag:
				hasIntFlag = true
			case *urfavecli.BoolFlag:
				hasBoolFlag = true
			}
		}

		assert.True(t, hasIntFlag, "Should have an IntFlag")
		assert.True(t, hasBoolFlag, "Should have a BoolFlag")
	})
}

// TestConvertCommandHandler_Integration tests command handler integration
func TestConvertCommandHandler_Integration(t *testing.T) {
	t.Run("BuildCommand_ActionIsCallable", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		assert.NotNil(t, cmd.Action)
	})

	t.Run("BuildCommand_ProducesValidUrfaveCommand", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()
		cmd := handler.BuildCommand()

		// Verify it's a valid urfave/cli Command
		assert.NotNil(t, cmd)
		assert.IsType(t, &urfavecli.Command{}, cmd)
	})

	t.Run("Handler_CanBeReused", func(t *testing.T) {
		handler := commands.NewConvertCommandHandler()

		cmd1 := handler.BuildCommand()
		cmd2 := handler.BuildCommand()

		assert.Equal(t, cmd1.Name, cmd2.Name)
		assert.Equal(t, len(cmd1.Flags), len(cmd2.Flags))
	})
}
