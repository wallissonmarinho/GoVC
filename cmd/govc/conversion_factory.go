package main

import (
	"fmt"
	"runtime"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/commands"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

// BuildConversionExecutorFactory returns a commands.ExecutorFactory which constructs
// the Conversion executor (wires discovery, converter, filesystem, reporter).
func BuildConversionExecutorFactory() commands.ExecutorFactory {
	return func(c *urfavecli.Context) (ports.Executor, error) {
		workers := c.Int("workers")
		saveLogs := c.Bool("logs")
		inputDir := c.Args().First()

		if inputDir == "" {
			return nil, fmt.Errorf("input directory is required")
		}

		if workers <= 0 {
			workers = runtime.NumCPU()
		}

		cliConfig, err := cli.NewCLIConfigFromContext(workers, saveLogs, inputDir)
		if err != nil {
			return nil, err
		}

		// Adapters: Discovery & Conversion
		discoveryAdapter := filesystem.NewFilesystemAdapter()
		converterAdapter := ffmpeg.NewFFmpegAdapter()
		fileSystemAdapter := filesystem.NewFilesystemAdapter()

		// Adapters: Output (Reporter)
		reporterAdapter := cli.NewLoggerReporter()

		// Services (Core Business Logic)
		conversionService := services.NewConversionService(
			discoveryAdapter,
			converterAdapter,
			fileSystemAdapter,
			reporterAdapter,
			cliConfig,
		)

		return conversionService, nil
	}
}
