package ffmpeg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wallissonmarinho/GoVC/internal/core/domain"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// Ensure FFmpegAdapter implements VideoConverterPort
var _ ports.VideoConverterPort = (*FFmpegAdapter)(nil)

// FFmpegAdapter is an output adapter for video conversion using ffmpeg.
type FFmpegAdapter struct{}

// NewFFmpegAdapter creates a new ffmpeg adapter.
func NewFFmpegAdapter() *FFmpegAdapter {
	return &FFmpegAdapter{}
}

// ConvertWithProgress converts a video to MP4 using ffmpeg with progress callback.
func (fa *FFmpegAdapter) ConvertWithProgress(video *domain.Video, inputDir string, progressCallback func(float64)) error {
	args := fa.buildFFmpegArgs(video, inputDir)

	cmd := exec.Command("ffmpeg", args...)

	// Open log file for stderr
	logFile, err := os.Create(video.LogPath())
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	cmd.Stderr = logFile

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Parse progress in a goroutine to allow parallelism
	go fa.parseProgress(stdout, video, progressCallback)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}

// buildFFmpegArgs builds the ffmpeg command arguments.
func (fa *FFmpegAdapter) buildFFmpegArgs(video *domain.Video, inputDir string) []string {
	baseArgs := []string{
		"-i", video.Path,
		"-map", "0:v",
		"-map", "0:a",
		"-map", "0:s?",
		"-c:v", "copy",
		"-c:a", "aac", "-b:a", "192k",
		"-c:s", "mov_text",
		"-progress", "pipe:1",
		"-nostats",
		video.OutputPath(),
	}

	// If external SRT exists, add it
	srtPath := video.SubtitlePath(inputDir)
	if _, err := os.Stat(srtPath); err == nil {
		return []string{
			"-i", video.Path,
			"-i", srtPath,
			"-map", "0:v",
			"-map", "0:a",
			"-map", "0:s?",
			"-map", "1:s?",
			"-c:v", "copy",
			"-c:a", "aac", "-b:a", "192k",
			"-c:s", "mov_text",
			"-progress", "pipe:1",
			"-nostats",
			video.OutputPath(),
		}
	}

	return baseArgs
}

// parseProgress parses ffmpeg progress output.
func (fa *FFmpegAdapter) parseProgress(stdout interface{}, video *domain.Video, progressCallback func(float64)) {
	reader, ok := stdout.(interface {
		Read(p []byte) (n int, err error)
	})
	if !ok {
		log.Println("Error: stdout does not support reading")
		return
	}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		switch key {
		case "out_time_ms":
			percent := fa.calculateProgressFromMicroseconds(value, video.Duration)
			if percent >= 0 {
				progressCallback(percent)
			}
		case "out_time":
			seconds, err := parseOutTime(value)
			if err == nil && video.Duration > 0 {
				percent := (seconds / video.Duration) * 100.0
				progressCallback(percent)
			}
		case "progress":
			if value == "end" {
				progressCallback(100)
			}
		}
	}
}

// calculateProgressFromMicroseconds converts microseconds to a percentage.
func (fa *FFmpegAdapter) calculateProgressFromMicroseconds(us string, duration float64) float64 {
	if duration <= 0 {
		return -1
	}

	microseconds, err := strconv.ParseFloat(us, 64)
	if err != nil {
		return -1
	}

	seconds := microseconds / 1000000.0
	return (seconds / duration) * 100.0
}

// GetDuration gets the duration of a video using ffprobe.
func (fa *FFmpegAdapter) GetDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0, fmt.Errorf("ffprobe returned empty duration")
	}

	duration, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

// HasExternalSubtitles checks if external SRT subtitles exist.
func (fa *FFmpegAdapter) HasExternalSubtitles(video *domain.Video, inputDir string) bool {
	srtPath := video.SubtitlePath(inputDir)
	_, err := os.Stat(srtPath)
	return err == nil
}

// parseOutTime parses HH:MM:SS.ms format to seconds.
func parseOutTime(s string) (float64, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid out_time format: %s", s)
	}

	h, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	m, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}

	sec, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, err
	}

	return h*3600 + m*60 + sec, nil
}
