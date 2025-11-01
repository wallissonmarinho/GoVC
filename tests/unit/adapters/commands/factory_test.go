package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoVC/internal/adapters/commands"
)

// TestCommandFactory tests the command factory
func TestCommandFactory(t *testing.T) {
	t.Run("NewCommandFactory", func(t *testing.T) {
		factory := commands.NewCommandFactory()

		assert.NotNil(t, factory)
	})

	t.Run("BuildCommands_ReturnsNonEmpty", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		assert.NotEmpty(t, cmds)
	})

	t.Run("BuildCommands_IncludesConvertCommand", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		// Find convert command
		found := false
		for _, cmd := range cmds {
			if cmd.Name == "convert" {
				found = true
				break
			}
		}

		assert.True(t, found, "convert command should be in the list")
	})

	t.Run("BuildCommands_ConvertCommandHasAction", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		// Find convert command
		for _, cmd := range cmds {
			if cmd.Name == "convert" {
				assert.NotNil(t, cmd.Action, "convert command should have an action")
				break
			}
		}
	})

	t.Run("BuildCommands_CallableTwice", func(t *testing.T) {
		factory := commands.NewCommandFactory()

		cmds1 := factory.BuildCommands()
		cmds2 := factory.BuildCommands()

		assert.Equal(t, len(cmds1), len(cmds2))
	})

	t.Run("BuildCommands_EachCommandHasName", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		for _, cmd := range cmds {
			assert.NotEmpty(t, cmd.Name)
		}
	})

	t.Run("BuildCommands_EachCommandHasUsage", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		for _, cmd := range cmds {
			assert.NotEmpty(t, cmd.Usage)
		}
	})
}

// TestCommandFactory_Integration tests factory integration
func TestCommandFactory_Integration(t *testing.T) {
	t.Run("Factory_ProducesValidCLICommands", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		assert.NotNil(t, cmds)
		assert.True(t, len(cmds) > 0)

		// Each command should be ready for urfave/cli
		for _, cmd := range cmds {
			assert.NotNil(t, cmd)
			assert.NotEmpty(t, cmd.Name)
			assert.NotNil(t, cmd.Action)
		}
	})

	t.Run("Factory_AllCommandsAreDistinct", func(t *testing.T) {
		factory := commands.NewCommandFactory()
		cmds := factory.BuildCommands()

		names := make(map[string]bool)
		for _, cmd := range cmds {
			assert.False(t, names[cmd.Name], "command name should be unique: "+cmd.Name)
			names[cmd.Name] = true
		}
	})
}
