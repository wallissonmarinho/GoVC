package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wallissonmarinho/GoVC/internal/core/domain"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// ConversionService is the application use case for converting videos.
// It depends only on interfaces (ports), not concrete implementations.
type ConversionService struct {
	discovery  ports.VideoDiscoveryPort
	converter  ports.VideoConverterPort
	fileSystem ports.FileSystemPort
	reporter   ports.ProgressReporterPort
	config     ports.ConfigPort
}

// NewConversionService creates a new conversion service.
func NewConversionService(
	discovery ports.VideoDiscoveryPort,
	converter ports.VideoConverterPort,
	fileSystem ports.FileSystemPort,
	reporter ports.ProgressReporterPort,
	config ports.ConfigPort,
) *ConversionService {
	return &ConversionService{
		discovery:  discovery,
		converter:  converter,
		fileSystem: fileSystem,
		reporter:   reporter,
		config:     config,
	}
}

// Execute runs the video conversion process.
func (cs *ConversionService) Execute() error {
	// Discover videos
	videos, err := cs.discovery.FindVideos(cs.config.GetInputDir())
	if err != nil {
		return fmt.Errorf("failed to discover videos: %w", err)
	}

	if len(videos) == 0 {
		log.Println("No .mkv files found to convert.")
		return nil
	}

	// Create output directory
	if err := cs.discovery.CreateOutputDir(cs.config.GetOutputDir()); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Update video paths and durations
	for _, v := range videos {
		v.OutputDir = cs.config.GetOutputDir()
		if duration, err := cs.converter.GetDuration(v.Path); err == nil {
			v.Duration = duration
		}
		v.HasSubs = cs.converter.HasExternalSubtitles(v, cs.config.GetInputDir())
	}

	// Setup progress tracking and printer
	tracker := domain.NewProgressTracker(len(videos))
	go cs.printProgressPeriodically(tracker)

	// Run conversions with worker pool
	semaphore := make(chan struct{}, cs.config.GetWorkers())
	var wg sync.WaitGroup

	for _, video := range videos {
		video := video // capture
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			cs.convertVideo(video, tracker)
		}()
	}

	wg.Wait()
	return nil
}

// convertVideo converts a single video and handles errors.
func (cs *ConversionService) convertVideo(video *domain.Video, tracker *domain.ProgressTracker) {
	defer tracker.MarkCompleted()

	cs.reporter.ReportConversionStart(video.Filename(), video.HasSubs)

	// Execute conversion with progress callback
	if err := cs.converter.ConvertWithProgress(video, cs.config.GetInputDir(), func(percent float64) {
		tracker.Update(video.BaseName, percent)
	}); err != nil {
		cs.reporter.ReportError(fmt.Sprintf("Failed to convert %s: %v", video.Filename(), err))
		return
	}

	// Validate output
	if !cs.fileSystem.IsValidOutput(video.OutputPath()) {
		cs.reporter.ReportError(fmt.Sprintf("Output file invalid for %s", video.Filename()))
		return
	}

	// Clean up log if NOT saving logs
	if !cs.config.SaveLogsEnabled() {
		_ = cs.fileSystem.RemoveFile(video.LogPath())
	}

	cs.reporter.ReportConversionFinish(video.Filename(), video.OutputPath(), true)
	tracker.Update(video.BaseName, 100)
}

// printProgressPeriodically prints progress at regular intervals.
func (cs *ConversionService) printProgressPeriodically(tracker *domain.ProgressTracker) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		snapshot := tracker.GetSnapshot()
		cs.reporter.ReportProgress(snapshot, tracker.IsComplete())
		if tracker.IsComplete() {
			return
		}
	}
}
