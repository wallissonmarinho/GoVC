package cli

import (
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// ConvertCommand implements ports.ServiceCommand
type ConvertCommand struct {
	executor ports.Executor
	name     string
}

var _ ports.ServiceCommand = (*ConvertCommand)(nil)

// NewConvertCommand creates a new ConvertCommand
func NewConvertCommand(executor ports.Executor, name string) ports.ServiceCommand {
	return &ConvertCommand{executor: executor, name: name}
}

// Execute runs the executor
func (cc *ConvertCommand) Execute() error {
	return cc.executor.Execute()
}

// Name returns the command name
func (cc *ConvertCommand) Name() string {
	return cc.name
}
