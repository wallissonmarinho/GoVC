package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

func printHelp() {
	help := `GoVC — MKV to MP4 Batch Converter

USAGE:
  govc [FLAGS] <input-directory>

FLAGS:
  -cmd string
      Command to execute (default: "convert")
  
  -p int
      Number of parallel workers (default: number of CPUs)
  
  -logs bool
      Save per-file logs (default: true)
      Use -logs=false to delete successful conversion logs
  
  -help
      Show this help message

EXAMPLES:
  govc /path/to/videos
      Convert all MKV files using default settings
  
  govc -cmd convert -p 4 /path/to/videos
      Convert with 4 parallel workers
  
  govc -cmd convert -p 2 -logs=false /path/to/videos
      Convert with 2 workers, delete successful logs

DOCUMENTATION:
  For more information, see:
  - README.md
  - HEXAGONAL_ARCHITECTURE.md
  - EXTENSION_GUIDE.md
`
	fmt.Println(help)
}

func main() {
	// Parse flags
	cmdFlag := flag.String("cmd", "convert", "Command to execute: convert")
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()

	// Show help
	if *helpFlag {
		printHelp()
		return
	}

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

	// Execute specified command
	if err := executor.Execute(*cmdFlag); err != nil {
		log.Fatalf("❌ Command failed: %v", err)
	}
}
