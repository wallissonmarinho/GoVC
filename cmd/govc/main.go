package main

import (
	"log"

	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

func main() {
	// Adapters: Input (left side of hexagon)
	cliConfig, err := cli.NewCLIConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Adapters: Output (right side of hexagon)
	discoveryAdapter := filesystem.NewFilesystemAdapter()
	converterAdapter := ffmpeg.NewFFmpegAdapter()
	fileSystemAdapter := filesystem.NewFilesystemAdapter()
	reporterAdapter := cli.NewLoggerReporter()

	// Core: Application service (use case)
	conversionService := services.NewConversionService(
		discoveryAdapter,
		converterAdapter,
		fileSystemAdapter,
		reporterAdapter,
		cliConfig,
	)

	// Execute
	if err := conversionService.Execute(); err != nil {
		log.Fatalf("Conversion failed: %v", err)
	}
}
