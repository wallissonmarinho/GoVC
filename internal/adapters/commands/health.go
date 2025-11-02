package commands

import (
	"fmt"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
)

// HealthCommandHandler handles a simple health check command execution
type HealthCommandHandler struct {
	executorFactory ExecutorFactory
}

// NewHealthCommandHandler creates a new health command handler
func NewHealthCommandHandler(factory ExecutorFactory) *HealthCommandHandler {
	return &HealthCommandHandler{executorFactory: factory}
}

// Execute runs the health check via the provided executor factory
func (h *HealthCommandHandler) Execute(c *urfavecli.Context) error {
	if h.executorFactory == nil {
		return urfavecli.Exit("executor factory not provided", 1)
	}

	exec, err := h.executorFactory(c)
	if err != nil {
		return urfavecli.Exit(fmt.Sprintf("Configuration error: %v", err), 1)
	}

	executor := cli.NewCommandExecutor()
	executor.Register("health", cli.NewConvertCommand(exec, "health"))

	if err := executor.Execute("health"); err != nil {
		return urfavecli.Exit(fmt.Sprintf("❌ Command failed: %v", err), 1)
	}

	return nil
}

// BuildCommand returns the urfave/cli Command configuration for health
func (h *HealthCommandHandler) BuildCommand() *urfavecli.Command {
	return &urfavecli.Command{
		Name:   "health",
		Usage:  "Run a quick health check (validate input/output dirs)",
		Action: h.Execute,
	}
}
