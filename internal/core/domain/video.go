package domain

import (
	"path/filepath"
	"strings"
)

// Video represents a media file to be converted.
type Video struct {
	Path      string
	BaseName  string
	HasSubs   bool
	Duration  float64
	OutputDir string
}

// NewVideo creates a new Video from a file path.
func NewVideo(filePath string, outputDir string) *Video {
	filename := filepath.Base(filePath)
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

	return &Video{
		Path:      filePath,
		BaseName:  baseName,
		Duration:  0,
		OutputDir: outputDir,
	}
}

// OutputPath returns the output MP4 file path.
func (v *Video) OutputPath() string {
	return filepath.Join(v.OutputDir, v.BaseName+".mp4")
}

// LogPath returns the log file path.
func (v *Video) LogPath() string {
	return filepath.Join(v.OutputDir, v.BaseName+".log")
}

// SubtitlePath returns the expected SRT subtitle file path.
func (v *Video) SubtitlePath(inputDir string) string {
	return filepath.Join(inputDir, v.BaseName+".srt")
}

// Filename returns just the filename with extension.
func (v *Video) Filename() string {
	return filepath.Base(v.Path)
}
