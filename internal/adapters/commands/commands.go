package commands

import (
	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// ExecutorFactory builds a ports.Executor given a urfave/cli Context.
// This lets the composition/wiring be provided by the application (or tests)
// while keeping the adapter decoupled from concrete implementations.
type ExecutorFactory func(c *urfavecli.Context) (ports.Executor, error)
