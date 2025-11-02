package main

import (
	"log"
	"os"

	urfavecli "github.com/urfave/cli/v2"
)

func main() {
	// Build commands using factory pattern with a provided executor factory
	commandFactory := BuildCommandFactory()

	app := &urfavecli.App{
		Name:        "govc",
		Usage:       "MKV to MP4 Batch Converter",
		Description: "Batch convert MKV videos to MP4 with parallel processing and subtitle support",
		Commands:    commandFactory.BuildCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("❌ Error: %v", err)
	}
}
