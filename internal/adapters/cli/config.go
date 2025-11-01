package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// CLIConfig is an input adapter that parses CLI flags.
type CLIConfig struct {
	inputDir  string
	outputDir string
	workers   int
	saveLogs  bool
}

// NewCLIConfig creates a new CLI configuration by parsing flags.
// Deprecated: Use NewCLIConfigFromContext instead.
func NewCLIConfig() (*CLIConfig, error) {
	workers := flag.Int("p", runtime.NumCPU(), "number of parallel ffmpeg processes")
	saveLogs := flag.Bool("logs", false, "save per-file ffmpeg logs to mp4/*.log")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: govc [-p <num_workers>] [-logs] <directory>\n  -logs: keep FFmpeg logs (by default they are deleted)")
	}

	dir := args[0]
	if *workers < 1 {
		*workers = 1
	}

	return &CLIConfig{
		inputDir:  dir,
		outputDir: filepath.Join(dir, "mp4"),
		workers:   *workers,
		saveLogs:  *saveLogs,
	}, nil
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
