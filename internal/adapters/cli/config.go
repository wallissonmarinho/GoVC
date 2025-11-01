package cli

import (
	"fmt"
	"path/filepath"

	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// CLIConfig is an input adapter that parses CLI flags.
type CLIConfig struct {
	inputDir  string
	outputDir string
	workers   int
	saveLogs  bool
}

// NewCLIConfigFromContext creates a new CLI configuration from urfave/cli context.
func NewCLIConfigFromContext(workers int, saveLogs bool, inputDir string) (*CLIConfig, error) {
	if inputDir == "" {
		return nil, fmt.Errorf("input directory is required")
	}

	if workers < 1 {
		workers = 1
	}

	return &CLIConfig{
		inputDir:  inputDir,
		outputDir: filepath.Join(inputDir, "mp4"),
		workers:   workers,
		saveLogs:  saveLogs,
	}, nil
}

// GetInputDir returns the input directory.
func (c *CLIConfig) GetInputDir() string {
	return c.inputDir
}

// GetOutputDir returns the output directory.
func (c *CLIConfig) GetOutputDir() string {
	return c.outputDir
}

// GetWorkers returns the number of parallel workers.
func (c *CLIConfig) GetWorkers() int {
	return c.workers
}

// SaveLogsEnabled returns whether to save logs.
func (c *CLIConfig) SaveLogsEnabled() bool {
	return c.saveLogs
}

// Ensure CLIConfig implements ConfigPort
var _ ports.ConfigPort = (*CLIConfig)(nil)
