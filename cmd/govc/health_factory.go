package main

import (
	"fmt"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/commands"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

// BuildHealthExecutorFactory returns an ExecutorFactory that constructs a
// HealthCheckService to be used by the health command.
func BuildHealthExecutorFactory() commands.ExecutorFactory {
	return func(c *urfavecli.Context) (ports.Executor, error) {
		inputDir := c.Args().First()

		if inputDir == "" {
			return nil, fmt.Errorf("input directory is required")
		}

		// Health check doesn't depend on workers or log-saving behavior.
		// Pass sensible defaults directly to the CLI config constructor.
		cliConfig, err := cli.NewCLIConfigFromContext(1, false, inputDir)
		if err != nil {
			return nil, err
		}

		// Adapters: Discovery & FileSystem
		discoveryAdapter := filesystem.NewFilesystemAdapter()
		fileSystemAdapter := filesystem.NewFilesystemAdapter()

		// Services (Core Business Logic)
		healthService := services.NewHealthCheckService(
			discoveryAdapter,
			fileSystemAdapter,
			cliConfig,
		)

		return healthService, nil
	}
}
