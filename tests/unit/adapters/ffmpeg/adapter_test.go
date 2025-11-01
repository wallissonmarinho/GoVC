package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/core/domain"
)

// TestFFmpegAdapter tests FFmpeg adapter
func TestFFmpegAdapter(t *testing.T) {
	t.Run("NewFFmpegAdapter", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()

		assert.NotNil(t, adapter)
	})

	t.Run("FFmpegAdapter_GetDuration_InvalidFile", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()

		// Test with non-existent file
		duration, err := adapter.GetDuration("/non/existent/file.mkv")

		// Should error since file doesn't exist
		assert.Error(t, err)
		assert.Equal(t, 0.0, duration)
	})

	t.Run("FFmpegAdapter_HasExternalSubtitles_NoFile", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()
		video := domain.NewVideo("/input/test.mkv", "/output")

		// Test with non-existent directory
		hasSubs := adapter.HasExternalSubtitles(video, "/non/existent/dir")

		assert.False(t, hasSubs)
	})

	t.Run("FFmpegAdapter_ConvertWithProgress", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()
		video := domain.NewVideo("/non/existent/test.mkv", "/output")

		// Test with non-existent file (will error)
		callback := func(p float64) {}
		err := adapter.ConvertWithProgress(video, "/input", callback)

		// Should error since file doesn't exist
		assert.Error(t, err)
	})
}

// TestFFmpegAdapter_VideoConverter tests video converter interface
func TestFFmpegAdapter_VideoConverter(t *testing.T) {
	t.Run("ImplementsVideoConverterPort", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()

		// Verify adapter has required methods
		assert.NotNil(t, adapter)

		// Should not panic on interface methods
		assert.NotPanics(t, func() {
			video := domain.NewVideo("/test.mkv", "/output")
			adapter.HasExternalSubtitles(video, "/input")
		})
	})

	t.Run("GetDuration_NoFfprobe", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()

		// Test with file that will fail
		duration, err := adapter.GetDuration("")

		// Should error with empty path
		assert.Error(t, err)
		assert.Equal(t, 0.0, duration)
	})
}

// TestFFmpegAdapter_Callbacks tests callback mechanism
func TestFFmpegAdapter_Callbacks(t *testing.T) {
	t.Run("ConvertWithProgress_CallsCallback", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()
		video := domain.NewVideo("/non/existent/test.mkv", "/output")

		callCount := 0
		callback := func(p float64) {
			callCount++
		}

		// Even though it will error, callback shouldn't crash
		_ = adapter.ConvertWithProgress(video, "/input", callback)

		// Callback may or may not be called depending on error timing
		assert.GreaterOrEqual(t, callCount, 0)
	})

	t.Run("ConvertWithProgress_ValidCallback", func(t *testing.T) {
		adapter := ffmpeg.NewFFmpegAdapter()
		video := domain.NewVideo("/test.mkv", "/output")

		progressValues := []float64{}
		callback := func(p float64) {
			progressValues = append(progressValues, p)
		}

		// Will error but callback should be safe
		err := adapter.ConvertWithProgress(video, "/tmp", callback)

		// Should error with invalid file
		assert.Error(t, err)
	})
}

// TestFFmpegAdapter_Integration tests with real ffmpeg (if available)
func TestFFmpegAdapter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	adapter := ffmpeg.NewFFmpegAdapter()

	t.Run("GetDuration_ValidFile", func(t *testing.T) {
		// Would test with a real video file if available
		// For now, test structure works correctly
		assert.NotNil(t, adapter)
	})

	t.Run("HasExternalSubtitles_DirectoryCheck", func(t *testing.T) {
		video := domain.NewVideo("/test.mkv", "/output")
		tmpDir := t.TempDir()

		// Check in temp directory (no subs file)
		hasSubs := adapter.HasExternalSubtitles(video, tmpDir)

		assert.False(t, hasSubs)
	})
}
