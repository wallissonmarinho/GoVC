package domain

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoVC/internal/core/domain"
)

// TestVideo is the main test suite for Video domain model
func TestVideo(t *testing.T) {
	t.Run("NewVideo", func(t *testing.T) {
		filePath := "/input/movie.mkv"
		outputDir := "/output"

		video := domain.NewVideo(filePath, outputDir)

		assert.NotNil(t, video)
		assert.Equal(t, filePath, video.Path)
		assert.Equal(t, "movie", video.BaseName)
		assert.Equal(t, outputDir, video.OutputDir)
		assert.Equal(t, 0.0, video.Duration)
		assert.False(t, video.HasSubs)
	})

	t.Run("NewVideoWithComplexPath", func(t *testing.T) {
		filePath := "/home/user/Videos/movies/my-movie-2024.mkv"
		outputDir := "/output"

		video := domain.NewVideo(filePath, outputDir)

		assert.Equal(t, "my-movie-2024", video.BaseName)
		assert.Equal(t, filePath, video.Path)
	})

	t.Run("OutputPath", func(t *testing.T) {
		video := domain.NewVideo("/input/movie.mkv", "/output")
		expected := filepath.Join("/output", "movie.mp4")

		assert.Equal(t, expected, video.OutputPath())
	})

	t.Run("LogPath", func(t *testing.T) {
		video := domain.NewVideo("/input/movie.mkv", "/output")
		expected := filepath.Join("/output", "movie.log")

		assert.Equal(t, expected, video.LogPath())
	})

	t.Run("SubtitlePath", func(t *testing.T) {
		video := domain.NewVideo("/input/movie.mkv", "/output")
		inputDir := "/input"
		expected := filepath.Join("/input", "movie.srt")

		assert.Equal(t, expected, video.SubtitlePath(inputDir))
	})

	t.Run("Filename", func(t *testing.T) {
		video := domain.NewVideo("/input/movie.mkv", "/output")

		assert.Equal(t, "movie.mkv", video.Filename())
	})

	t.Run("FilenameWithComplexPath", func(t *testing.T) {
		video := domain.NewVideo("/home/user/Videos/films/my-movie-2024.mkv", "/output")

		assert.Equal(t, "my-movie-2024.mkv", video.Filename())
	})

	t.Run("VideoPathsConsistency", func(t *testing.T) {
		video := domain.NewVideo("/input/movie.mkv", "/output")

		outputPath := video.OutputPath()
		logPath := video.LogPath()
		subtitlePath := video.SubtitlePath("/input")

		// All should contain the BaseName
		assert.Contains(t, outputPath, video.BaseName)
		assert.Contains(t, logPath, video.BaseName)
		assert.Contains(t, subtitlePath, video.BaseName)
	})

	t.Run("VideoWithMultipleExtensions", func(t *testing.T) {
		tests := []struct {
			name         string
			filePath     string
			expectedBase string
		}{
			{
				name:         "mkv extension",
				filePath:     "/input/movie.mkv",
				expectedBase: "movie",
			},
			{
				name:         "complex filename with dots",
				filePath:     "/input/movie.2024.mkv",
				expectedBase: "movie.2024",
			},
			{
				name:         "filename with multiple dots and mkv",
				filePath:     "/input/movie.1080p.5.1.mkv",
				expectedBase: "movie.1080p.5.1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				video := domain.NewVideo(tt.filePath, "/output")
				assert.Equal(t, tt.expectedBase, video.BaseName)
			})
		}
	})

	t.Run("VideoModification", func(t *testing.T) {
		video := domain.NewVideo("/input/movie.mkv", "/output")

		// Modify fields
		video.Duration = 120.5
		video.HasSubs = true

		assert.Equal(t, 120.5, video.Duration)
		assert.True(t, video.HasSubs)
	})
}
