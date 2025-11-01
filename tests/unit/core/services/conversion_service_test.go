package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/ffmpeg"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/domain"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

// TestConversionService is the main test suite for ConversionService
func TestConversionService(t *testing.T) {
	t.Run("NewConversionService", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		assert.NotNil(t, service)
	})

	t.Run("ExecuteNoVideosFound", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		mockConfig.On("GetInputDir").Return("/input")
		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{}, nil)

		err := service.Execute()

		assert.NoError(t, err)
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
		mockConverter.AssertExpectations(t)
		mockFileSystem.AssertExpectations(t)
		mockReporter.AssertExpectations(t)
	})

	t.Run("ExecuteDiscoveryError", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		mockConfig.On("GetInputDir").Return("/input")
		mockDiscovery.On("FindVideos", "/input").Return(nil, errors.New("discovery failed"))

		err := service.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to discover videos")
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteCreateOutputDirError", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video := domain.NewVideo("/input/test.mkv", "/output")
		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(errors.New("mkdir failed"))

		err := service.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create output directory")
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteSingleVideoConversionSuccess", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video := domain.NewVideo("/input/test.mkv", "/output")

		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockConfig.On("GetWorkers").Return(1)
		mockConfig.On("SaveLogsEnabled").Return(false)

		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(nil)

		mockConverter.On("GetDuration", video.Path).Return(100.5, nil)
		mockConverter.On("HasExternalSubtitles", video, "/input").Return(false)
		mockConverter.On("ConvertWithProgress", video, "/input", mock.Anything).Return(nil)

		mockFileSystem.On("IsValidOutput", video.OutputPath()).Return(true)
		mockFileSystem.On("RemoveFile", video.LogPath()).Return(nil)

		mockReporter.On("ReportConversionStart", video.Filename(), false).Return()
		mockReporter.On("ReportConversionFinish", video.Filename(), video.OutputPath(), true).Return()
		mockReporter.On("ReportProgress", mock.Anything, mock.Anything).Return().Maybe()

		err := service.Execute()

		assert.NoError(t, err)
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteConversionFailure", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video := domain.NewVideo("/input/test.mkv", "/output")

		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockConfig.On("GetWorkers").Return(1)

		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(nil)

		mockConverter.On("GetDuration", video.Path).Return(100.5, nil)
		mockConverter.On("HasExternalSubtitles", video, "/input").Return(false)
		mockConverter.On("ConvertWithProgress", video, "/input", mock.Anything).Return(errors.New("ffmpeg error"))

		mockReporter.On("ReportConversionStart", video.Filename(), false).Return()
		mockReporter.On("ReportError", mock.MatchedBy(func(msg string) bool {
			return fmt.Sprintf("Failed to convert %s: ffmpeg error", video.Filename()) == msg
		})).Return()
		mockReporter.On("ReportProgress", mock.Anything, mock.Anything).Return().Maybe()

		err := service.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conversion completed with 1 error(s)")
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteInvalidOutput", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video := domain.NewVideo("/input/test.mkv", "/output")

		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockConfig.On("GetWorkers").Return(1)

		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(nil)

		mockConverter.On("GetDuration", video.Path).Return(100.5, nil)
		mockConverter.On("HasExternalSubtitles", video, "/input").Return(false)
		mockConverter.On("ConvertWithProgress", video, "/input", mock.Anything).Return(nil)

		mockFileSystem.On("IsValidOutput", video.OutputPath()).Return(false)

		mockReporter.On("ReportConversionStart", video.Filename(), false).Return()
		mockReporter.On("ReportError", mock.MatchedBy(func(msg string) bool {
			return fmt.Sprintf("Output file invalid for %s", video.Filename()) == msg
		})).Return()
		mockReporter.On("ReportProgress", mock.Anything, mock.Anything).Return().Maybe()

		err := service.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conversion completed with 1 error(s)")
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteMultipleVideos", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video1 := domain.NewVideo("/input/test1.mkv", "/output")
		video2 := domain.NewVideo("/input/test2.mkv", "/output")

		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockConfig.On("GetWorkers").Return(2)
		mockConfig.On("SaveLogsEnabled").Return(false)

		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video1, video2}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(nil)

		mockConverter.On("GetDuration", mock.Anything).Return(100.5, nil)
		mockConverter.On("HasExternalSubtitles", mock.Anything, "/input").Return(false)
		mockConverter.On("ConvertWithProgress", mock.Anything, "/input", mock.Anything).Return(nil)

		mockFileSystem.On("IsValidOutput", mock.Anything).Return(true)
		mockFileSystem.On("RemoveFile", mock.Anything).Return(nil)

		mockReporter.On("ReportConversionStart", mock.Anything, false).Return()
		mockReporter.On("ReportConversionFinish", mock.Anything, mock.Anything, true).Return()
		mockReporter.On("ReportProgress", mock.Anything, mock.Anything).Return().Maybe()

		err := service.Execute()

		assert.NoError(t, err)
		mockConverter.AssertNumberOfCalls(t, "ConvertWithProgress", 2)
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteWithLogsEnabled", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video := domain.NewVideo("/input/test.mkv", "/output")

		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockConfig.On("GetWorkers").Return(1)
		mockConfig.On("SaveLogsEnabled").Return(true)

		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(nil)

		mockConverter.On("GetDuration", video.Path).Return(100.5, nil)
		mockConverter.On("HasExternalSubtitles", video, "/input").Return(false)
		mockConverter.On("ConvertWithProgress", video, "/input", mock.Anything).Return(nil)

		mockFileSystem.On("IsValidOutput", video.OutputPath()).Return(true)

		mockReporter.On("ReportConversionStart", video.Filename(), false).Return()
		mockReporter.On("ReportConversionFinish", video.Filename(), video.OutputPath(), true).Return()
		mockReporter.On("ReportProgress", mock.Anything, mock.Anything).Return().Maybe()

		err := service.Execute()

		assert.NoError(t, err)
		mockFileSystem.AssertNotCalled(t, "RemoveFile")
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})

	t.Run("ExecuteWithExternalSubtitles", func(t *testing.T) {
		mockDiscovery := new(filesystem.MockFilesystemAdapter)
		mockConverter := new(ffmpeg.MockFFmpegAdapter)
		mockFileSystem := new(filesystem.MockFilesystemAdapter)
		mockReporter := new(cli.MockLoggerReporter)
		mockConfig := new(cli.MockCLIConfig)

		service := services.NewConversionService(
			mockDiscovery,
			mockConverter,
			mockFileSystem,
			mockReporter,
			mockConfig,
		)

		video := domain.NewVideo("/input/test.mkv", "/output")

		mockConfig.On("GetInputDir").Return("/input")
		mockConfig.On("GetOutputDir").Return("/output")
		mockConfig.On("GetWorkers").Return(1)
		mockConfig.On("SaveLogsEnabled").Return(false)

		mockDiscovery.On("FindVideos", "/input").Return([]*domain.Video{video}, nil)
		mockDiscovery.On("CreateOutputDir", "/output").Return(nil)

		mockConverter.On("GetDuration", video.Path).Return(100.5, nil)
		mockConverter.On("HasExternalSubtitles", video, "/input").Return(true)
		mockConverter.On("ConvertWithProgress", video, "/input", mock.Anything).Return(nil)

		mockFileSystem.On("IsValidOutput", video.OutputPath()).Return(true)
		mockFileSystem.On("RemoveFile", video.LogPath()).Return(nil)

		mockReporter.On("ReportConversionStart", video.Filename(), true).Return()
		mockReporter.On("ReportConversionFinish", video.Filename(), video.OutputPath(), true).Return()
		mockReporter.On("ReportProgress", mock.Anything, mock.Anything).Return().Maybe()

		err := service.Execute()

		assert.NoError(t, err)
		mockReporter.AssertCalled(t, "ReportConversionStart", video.Filename(), true)
		mockConfig.AssertExpectations(t)
		mockDiscovery.AssertExpectations(t)
	})
}
