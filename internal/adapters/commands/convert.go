package commands

import (
	"fmt"
	"runtime"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

// ConvertCommandHandler handles the convert command execution
type ConvertCommandHandler struct{}

// NewConvertCommandHandler creates a new convert command handler
func NewConvertCommandHandler() *ConvertCommandHandler {
	return &ConvertCommandHandler{}
}

// Execute processes the convert command with urfave/cli context
func (h *ConvertCommandHandler) Execute(c *urfavecli.Context) error {
	// Extract arguments from urfave/cli context
	workers := c.Int("workers")
	saveLogs := c.Bool("logs")
	inputDir := c.Args().First()

	if inputDir == "" {
		return urfavecli.Exit("❌ Input directory is required", 1)
	}

	// Apply default workers if not specified or invalid
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	// Adapters: Input
	cliConfig, err := cli.NewCLIConfigFromContext(workers, saveLogs, inputDir)
	if err != nil {
		return urfavecli.Exit(fmt.Sprintf("Configuration error: %v", err), 1)
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

	// Commands (Adapters) - execute conversion
	executor := cli.NewCommandExecutor()
	executor.Register("convert", cli.NewConvertCommand(conversionService, "conversion"))

	// Execute the command
	if err := executor.Execute("convert"); err != nil {
		return urfavecli.Exit(fmt.Sprintf("❌ Command failed: %v", err), 1)
	}

	return nil
}

// BuildCommand returns the urfave/cli Command configuration for convert
func (h *ConvertCommandHandler) BuildCommand() *urfavecli.Command {
	return &urfavecli.Command{
		Name:  "convert",
		Usage: "Convert MKV files to MP4",
		Description: "Batch convert MKV videos to MP4 format with parallel processing.\n" +
			"Supports subtitle embedding and FFmpeg logging.",
		Flags: []urfavecli.Flag{
			&urfavecli.IntFlag{
				Name:    "workers",
				Aliases: []string{"p"},
				Usage:   "Number of parallel workers (default: number of CPUs)",
			},
			&urfavecli.BoolFlag{
				Name:    "logs",
				Aliases: []string{"l"},
				Value:   false,
				Usage:   "Keep per-file logs (default: delete successful logs)",
			},
		},
		Action: h.Execute,
	}
}
