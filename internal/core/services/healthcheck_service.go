package services

import (
	"fmt"

	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// HealthCheckService is a small example service that validates input/output
// directories and ensures the output directory exists. It implements
// ports.Executor so it can be used as a command executor in the CLI.
type HealthCheckService struct {
	discovery  ports.VideoDiscoveryPort
	fileSystem ports.FileSystemPort
	config     ports.ConfigPort
}

// NewHealthCheckService creates a new HealthCheckService.
func NewHealthCheckService(discovery ports.VideoDiscoveryPort, fs ports.FileSystemPort, cfg ports.ConfigPort) *HealthCheckService {
	return &HealthCheckService{
		discovery:  discovery,
		fileSystem: fs,
		config:     cfg,
	}
}

// Execute runs the health check. It verifies the input directory exists and
// attempts to create the output directory if necessary.
func (h *HealthCheckService) Execute() error {
	input := h.config.GetInputDir()
	output := h.config.GetOutputDir()

	if !h.fileSystem.FileExists(input) {
		return fmt.Errorf("input directory not found: %s", input)
	}

	// Ensure output directory exists
	if err := h.discovery.CreateOutputDir(output); err != nil {
		return fmt.Errorf("failed to ensure output dir: %w", err)
	}

	return nil
}
