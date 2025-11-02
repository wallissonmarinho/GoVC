package cli

import (
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// HealthCommand implements ports.ServiceCommand
type HealthCommand struct {
	executor ports.Executor
	name     string
}

var _ ports.ServiceCommand = (*HealthCommand)(nil)

// NewHealthCommand creates a new HealthCommand
func NewHealthCommand(executor ports.Executor, name string) ports.ServiceCommand {
	return &HealthCommand{executor: executor, name: name}
}

// Execute runs the executor
func (hc *HealthCommand) Execute() error {
	return hc.executor.Execute()
}

// Name returns the command name
func (hc *HealthCommand) Name() string {
	return hc.name
}
