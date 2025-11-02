package main

import "github.com/wallissonmarinho/GoVC/internal/adapters/commands"

// BuildCommandFactory composes and returns a CommandFactory wired with the
// executor factory. Keeps composition in one place for easier maintenance.
func BuildCommandFactory() *commands.CommandFactory {
	return commands.NewCommandFactoryWithExecutorFactories(
		BuildConversionExecutorFactory(),
		BuildHealthExecutorFactory(),
	)
}
