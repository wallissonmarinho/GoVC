package domain

import "sync"

// Progress represents the progress of a conversion.
type Progress struct {
	Filename string
	Percent  float64
}

// ProgressTracker tracks conversion progress for all videos.
type ProgressTracker struct {
	mu        sync.Mutex
	Progress  map[string]float64
	Total     int
	Completed int
}

// NewProgressTracker creates a new progress tracker.
func NewProgressTracker(total int) *ProgressTracker {
	return &ProgressTracker{
		Progress:  make(map[string]float64),
		Total:     total,
		Completed: 0,
	}
}

// Update updates the progress for a video.
func (pt *ProgressTracker) Update(filename string, percent float64) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if percent > 100 {
		percent = 100
	}
	pt.Progress[filename] = percent
}

// MarkCompleted increments the completed counter.
func (pt *ProgressTracker) MarkCompleted() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.Completed++
}

// IsComplete returns true if all conversions are done.
func (pt *ProgressTracker) IsComplete() bool {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	return pt.Completed >= pt.Total
}

// GetSnapshot returns a copy of current progress.
func (pt *ProgressTracker) GetSnapshot() map[string]float64 {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	snapshot := make(map[string]float64)
	for k, v := range pt.Progress {
		snapshot[k] = v
	}
	return snapshot
}
