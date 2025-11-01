package commands

import urfavecli "github.com/urfave/cli/v2"

// CommandFactory builds all available CLI commands
type CommandFactory struct {
	convertHandler *ConvertCommandHandler
}

// NewCommandFactory creates a new command factory
func NewCommandFactory() *CommandFactory {
	return &CommandFactory{
		convertHandler: NewConvertCommandHandler(),
	}
}

// BuildCommands returns all available CLI commands
func (f *CommandFactory) BuildCommands() []*urfavecli.Command {
	return []*urfavecli.Command{
		f.convertHandler.BuildCommand(),
		// Future commands can be added here:
		// f.validateHandler.BuildCommand(),
		// f.cleanupHandler.BuildCommand(),
	}
}
