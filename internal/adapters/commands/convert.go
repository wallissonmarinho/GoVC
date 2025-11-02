package commands

import (
	"fmt"

	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
)

// ConvertCommandHandler handles the convert command execution
type ConvertCommandHandler struct {
	executorFactory ExecutorFactory
}

// NewConvertCommandHandler creates a new convert command handler that uses
// the provided ExecutorFactory to obtain the concrete executor at runtime.
func NewConvertCommandHandler(factory ExecutorFactory) *ConvertCommandHandler {
	return &ConvertCommandHandler{executorFactory: factory}
}

// Execute processes the convert command with urfave/cli context
func (h *ConvertCommandHandler) Execute(c *urfavecli.Context) error {
	if h.executorFactory == nil {
		return urfavecli.Exit("executor factory not provided", 1)
	}

	exec, err := h.executorFactory(c)
	if err != nil {
		return urfavecli.Exit(fmt.Sprintf("Configuration error: %v", err), 1)
	}

	// Commands (Adapters) - execute conversion
	executor := cli.NewCommandExecutor()
	executor.Register("convert", cli.NewConvertCommand(exec, "conversion"))

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
