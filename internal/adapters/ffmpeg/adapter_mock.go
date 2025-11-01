package ffmpeg

import (
	"github.com/stretchr/testify/mock"
	"github.com/wallissonmarinho/GoVC/internal/core/domain"
)

// MockFFmpegAdapter is a mock implementation of VideoConverterPort for testing
type MockFFmpegAdapter struct {
	mock.Mock
}

func (m *MockFFmpegAdapter) ConvertWithProgress(video *domain.Video, inputDir string, progressCallback func(float64)) error {
	args := m.Called(video, inputDir, progressCallback)
	return args.Error(0)
}

func (m *MockFFmpegAdapter) GetDuration(videoPath string) (float64, error) {
	args := m.Called(videoPath)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockFFmpegAdapter) HasExternalSubtitles(video *domain.Video, inputDir string) bool {
	args := m.Called(video, inputDir)
	return args.Bool(0)
}
