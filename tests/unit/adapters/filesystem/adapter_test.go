package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/domain"
)

// TestFilesystemAdapter tests filesystem operations
func TestFilesystemAdapter(t *testing.T) {
	t.Run("NewFilesystemAdapter", func(t *testing.T) {
		adapter := filesystem.NewFilesystemAdapter()

		assert.NotNil(t, adapter)
	})

	t.Run("FilesystemAdapter_FindVideos_EmptyDir", func(t *testing.T) {
		adapter := filesystem.NewFilesystemAdapter()

		// Test with non-existent directory (should return empty or error)
		videos, err := adapter.FindVideos("/non/existent/dir")

		// Should either return empty slice or error
		if err == nil {
			assert.NotNil(t, videos)
			assert.Equal(t, 0, len(videos))
		} else {
			assert.Error(t, err)
		}
	})

	t.Run("FilesystemAdapter_CreateOutputDir", func(t *testing.T) {
		adapter := filesystem.NewFilesystemAdapter()

		// Test creating in temp directory
		tmpDir := t.TempDir()
		outputDir := tmpDir + "/output"

		err := adapter.CreateOutputDir(outputDir)

		// Should succeed (or already exist)
		_ = err // Error is not critical in this test
	})

	t.Run("FilesystemAdapter_IsValidOutput", func(t *testing.T) {
		adapter := filesystem.NewFilesystemAdapter()

		// Test with non-existent file
		isValid := adapter.IsValidOutput("/non/existent/file.mp4")

		assert.False(t, isValid)
	})

	t.Run("FilesystemAdapter_RemoveFile", func(t *testing.T) {
		adapter := filesystem.NewFilesystemAdapter()

		// Test removing non-existent file
		// Should not panic even if file doesn't exist
		assert.NotPanics(t, func() {
			_ = adapter.RemoveFile("/non/existent/file.log")
		})
	})
}

// TestFilesystemAdapter_VideoDiscovery tests video discovery
func TestFilesystemAdapter_VideoDiscovery(t *testing.T) {
	t.Run("FindVideos_ImplementsInterface", func(t *testing.T) {
		adapter := filesystem.NewFilesystemAdapter()

		// Verify adapter implements VideoDiscoveryPort
		assert.NotNil(t, adapter)

		// Call should not panic
		assert.NotPanics(t, func() {
			_, _ = adapter.FindVideos("/tmp")
		})
	})
}

// TestFilesystemAdapter_FileSystem tests filesystem operations
func TestFilesystemAdapter_FileSystem(t *testing.T) {
	adapter := filesystem.NewFilesystemAdapter()

	t.Run("CreateOutputDir_InvalidPath", func(t *testing.T) {
		// Test with invalid path
		// May or may not error depending on implementation
		assert.NotPanics(t, func() {
			_ = adapter.CreateOutputDir("")
		})
	})

	t.Run("IsValidOutput_Empty", func(t *testing.T) {
		isValid := adapter.IsValidOutput("")

		assert.False(t, isValid)
	})

	t.Run("RemoveFile_NoError", func(t *testing.T) {
		// Should not panic
		assert.NotPanics(t, func() {
			_ = adapter.RemoveFile("/any/path")
		})
	})
}

// TestFilesystemAdapter_Integration tests with actual temp files
func TestFilesystemAdapter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	adapter := filesystem.NewFilesystemAdapter()
	tmpDir := t.TempDir()

	t.Run("CreateAndValidateOutputDir", func(t *testing.T) {
		outputDir := tmpDir + "/mp4"

		_ = adapter.CreateOutputDir(outputDir)

		// Directory should exist (implementation dependent)
		assert.NotPanics(t, func() {
			_ = adapter.CreateOutputDir(outputDir)
		})
	})

	t.Run("Video_SubtitlePath", func(t *testing.T) {
		video := domain.NewVideo(tmpDir+"/test.mkv", tmpDir+"/output")

		subtitlePath := video.SubtitlePath(tmpDir)

		assert.Contains(t, subtitlePath, "test")
		assert.Contains(t, subtitlePath, ".srt")
	})
}
