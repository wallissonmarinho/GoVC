package main

import (
	"log"

	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

func main() {
	// Adapters: Input
	cliConfig, err := cli.NewCLIConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Adapters: Output
	discoveryAdapter := filesystem.NewFilesystemAdapter()
	converterAdapter := ffmpeg.NewFFmpegAdapter()
	fileSystemAdapter := filesystem.NewFilesystemAdapter()
	reporterAdapter := cli.NewLoggerReporter()

	// Services (Core)
	conversionService := services.NewConversionService(
		discoveryAdapter, converterAdapter, fileSystemAdapter, reporterAdapter, cliConfig,
	)

	// Commands (Adapters)
	executor := cli.NewCommandExecutor()
	executor.Register("convert", cli.NewConvertCommand(conversionService, "conversion"))

	// Execute default command
	_ = executor.Execute("convert")
}
