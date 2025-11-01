package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProgressTracker is the main test suite for ProgressTracker
func TestProgressTracker(t *testing.T) {
	t.Run("NewProgressTracker", func(t *testing.T) {
		total := 5
		tracker := NewProgressTracker(total)

		assert.NotNil(t, tracker)
		assert.Equal(t, total, tracker.Total)
		assert.Equal(t, 0, tracker.Completed)
		assert.Equal(t, 0, len(tracker.Progress))
	})

	t.Run("ProgressTrackerUpdate", func(t *testing.T) {
		tracker := NewProgressTracker(3)

		tracker.Update("video1.mkv", 50.0)
		tracker.Update("video2.mkv", 75.5)

		assert.Equal(t, 50.0, tracker.Progress["video1.mkv"])
		assert.Equal(t, 75.5, tracker.Progress["video2.mkv"])
	})

	t.Run("ProgressTrackerUpdateOverflow", func(t *testing.T) {
		tracker := NewProgressTracker(1)

		tracker.Update("video.mkv", 150.0) // Over 100%

		assert.Equal(t, 100.0, tracker.Progress["video.mkv"])
	})

	t.Run("ProgressTrackerUpdateOverwrite", func(t *testing.T) {
		tracker := NewProgressTracker(1)

		tracker.Update("video.mkv", 25.0)
		assert.Equal(t, 25.0, tracker.Progress["video.mkv"])

		tracker.Update("video.mkv", 75.0)
		assert.Equal(t, 75.0, tracker.Progress["video.mkv"])
	})

	t.Run("ProgressTrackerMarkCompleted", func(t *testing.T) {
		tracker := NewProgressTracker(3)

		tracker.MarkCompleted()
		assert.Equal(t, 1, tracker.Completed)

		tracker.MarkCompleted()
		assert.Equal(t, 2, tracker.Completed)

		tracker.MarkCompleted()
		assert.Equal(t, 3, tracker.Completed)
	})

	t.Run("ProgressTrackerIsComplete", func(t *testing.T) {
		tracker := NewProgressTracker(2)

		assert.False(t, tracker.IsComplete())

		tracker.MarkCompleted()
		assert.False(t, tracker.IsComplete())

		tracker.MarkCompleted()
		assert.True(t, tracker.IsComplete())
	})

	t.Run("ProgressTrackerIsCompleteExceeds", func(t *testing.T) {
		tracker := NewProgressTracker(2)

		tracker.MarkCompleted()
		tracker.MarkCompleted()
		tracker.MarkCompleted() // One more than total

		assert.True(t, tracker.IsComplete())
	})

	t.Run("ProgressTrackerGetSnapshot", func(t *testing.T) {
		tracker := NewProgressTracker(3)

		tracker.Update("video1.mkv", 30.0)
		tracker.Update("video2.mkv", 60.0)
		tracker.Update("video3.mkv", 90.0)

		snapshot := tracker.GetSnapshot()

		assert.Equal(t, 3, len(snapshot))
		assert.Equal(t, 30.0, snapshot["video1.mkv"])
		assert.Equal(t, 60.0, snapshot["video2.mkv"])
		assert.Equal(t, 90.0, snapshot["video3.mkv"])
	})

	t.Run("ProgressTrackerGetSnapshotIsolation", func(t *testing.T) {
		tracker := NewProgressTracker(1)

		tracker.Update("video.mkv", 50.0)
		snapshot1 := tracker.GetSnapshot()

		// Modify tracker after snapshot
		tracker.Update("video.mkv", 100.0)
		snapshot2 := tracker.GetSnapshot()

		// First snapshot should not be affected
		assert.Equal(t, 50.0, snapshot1["video.mkv"])
		assert.Equal(t, 100.0, snapshot2["video.mkv"])
	})

	t.Run("ProgressTrackerMultipleVideos", func(t *testing.T) {
		tracker := NewProgressTracker(5)

		// Simulate progress updates
		for i := 1; i <= 5; i++ {
			filename := fmt.Sprintf("video%d.mkv", i)
			tracker.Update(filename, float64(i)*20)
		}

		snapshot := tracker.GetSnapshot()
		assert.Equal(t, 5, len(snapshot))

		for i := 1; i <= 5; i++ {
			filename := fmt.Sprintf("video%d.mkv", i)
			assert.Equal(t, float64(i)*20, snapshot[filename])
		}
	})

	t.Run("ProgressTrackerZeroTotal", func(t *testing.T) {
		tracker := NewProgressTracker(0)

		assert.True(t, tracker.IsComplete()) // 0 >= 0 is true
	})

	t.Run("ProgressTrackerEmptySnapshot", func(t *testing.T) {
		tracker := NewProgressTracker(1)

		snapshot := tracker.GetSnapshot()

		assert.Equal(t, 0, len(snapshot))
	})

	t.Run("ProgressTrackerProgressUpdate", func(t *testing.T) {
		progress := Progress{
			Filename: "test.mkv",
			Percent:  75.5,
		}

		assert.Equal(t, "test.mkv", progress.Filename)
		assert.Equal(t, 75.5, progress.Percent)
	})
}
