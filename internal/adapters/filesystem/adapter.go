package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wallissonmarinho/GoVC/internal/core/domain"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// FilesystemAdapter is an output adapter for file system operations.
type FilesystemAdapter struct {
	outputDir string
}

// NewFilesystemAdapter creates a new filesystem adapter.
func NewFilesystemAdapter() *FilesystemAdapter {
	return &FilesystemAdapter{}
}

// FindVideos discovers all .mkv files in the given directory.
func (fa *FilesystemAdapter) FindVideos(dir string) ([]*domain.Video, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("input directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("input path is not a directory: %s", dir)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.mkv"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob MKV files: %w", err)
	}

	videos := make([]*domain.Video, 0, len(files))
	for _, f := range files {
		videos = append(videos, domain.NewVideo(f, ""))
	}

	return videos, nil
}

// CreateOutputDir creates the output directory if it doesn't exist.
func (fa *FilesystemAdapter) CreateOutputDir(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	return nil
}

// FileExists checks if a file exists.
func (fa *FilesystemAdapter) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsValidOutput checks if the output file exists and has non-zero size.
func (fa *FilesystemAdapter) IsValidOutput(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Size() > 0
}

// RemoveFile removes a file.
func (fa *FilesystemAdapter) RemoveFile(path string) error {
	return os.Remove(path)
}

// WriteLog writes lines to a log file.
func (fa *FilesystemAdapter) WriteLog(logPath string, lines []string) error {
	file, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// Ensure FilesystemAdapter implements VideoDiscoveryPort and FileSystemPort
var _ ports.VideoDiscoveryPort = (*FilesystemAdapter)(nil)
var _ ports.FileSystemPort = (*FilesystemAdapter)(nil)
