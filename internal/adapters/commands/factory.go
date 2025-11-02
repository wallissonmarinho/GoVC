package commands

import (
urfavecli "github.com/urfave/cli/v2"
"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// CommandFactory builds all available CLI commands
type CommandFactory struct {
	convertHandler *ConvertCommandHandler
	healthHandler  *HealthCommandHandler
}

// NewCommandFactoryWithExecutorFactories creates a CommandFactory with separate
// executor factories for conversion and health commands.
func NewCommandFactoryWithExecutorFactories(convFactory ExecutorFactory, healthFactory ExecutorFactory) *CommandFactory {
	return &CommandFactory{
		convertHandler: NewConvertCommandHandler(convFactory),
		healthHandler:  NewHealthCommandHandler(healthFactory),
	}
}

// NewCommandFactory is a convenience constructor used in tests or simple setups.
// It provides a no-op executor factory so commands exist but do not execute anything.
func NewCommandFactory() *CommandFactory {
	dummy := func(c *urfavecli.Context) (ports.Executor, error) { return nil, nil }
	return NewCommandFactoryWithExecutorFactories(dummy, dummy)
}

// BuildCommands returns all available CLI commands
func (f *CommandFactory) BuildCommands() []*urfavecli.Command {
	cmds := []*urfavecli.Command{}

	if f.convertHandler != nil {
		cmds = append(cmds, f.convertHandler.BuildCommand())
	}

	if f.healthHandler != nil {
		cmds = append(cmds, f.healthHandler.BuildCommand())
	}

	return cmds
}
