package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoVC/internal/adapters/cli"
	"github.com/wallissonmarinho/GoVC/internal/adapters/filesystem"
	"github.com/wallissonmarinho/GoVC/internal/core/services"
)

func TestHealthCheckService_Execute_Success(t *testing.T) {
	mockDiscovery := new(filesystem.MockFilesystemAdapter)
	mockFS := mockDiscovery // same mock implements both interfaces
	mockConfig := new(cli.MockCLIConfig)

	inputDir := "/some/input"
	outputDir := "/some/input/mp4"

	mockConfig.On("GetInputDir").Return(inputDir)
	mockConfig.On("GetOutputDir").Return(outputDir)

	mockFS.On("FileExists", inputDir).Return(true)
	mockDiscovery.On("CreateOutputDir", outputDir).Return(nil)

	svc := services.NewHealthCheckService(mockDiscovery, mockFS, mockConfig)

	err := svc.Execute()

	assert.NoError(t, err)
	mockConfig.AssertExpectations(t)
	mockFS.AssertExpectations(t)
	mockDiscovery.AssertExpectations(t)
}

func TestHealthCheckService_Execute_InputMissing(t *testing.T) {
	mockDiscovery := new(filesystem.MockFilesystemAdapter)
	mockFS := mockDiscovery
	mockConfig := new(cli.MockCLIConfig)

	inputDir := "/does/not/exist"
	outputDir := "/does/not/exist/mp4"

	mockConfig.On("GetInputDir").Return(inputDir)
	mockConfig.On("GetOutputDir").Return(outputDir)

	mockFS.On("FileExists", inputDir).Return(false)

	svc := services.NewHealthCheckService(mockDiscovery, mockFS, mockConfig)

	err := svc.Execute()

	assert.Error(t, err)
	mockConfig.AssertExpectations(t)
	mockFS.AssertExpectations(t)
}

func TestHealthCheckService_Execute_CreateOutputDirFails(t *testing.T) {
	mockDiscovery := new(filesystem.MockFilesystemAdapter)
	mockFS := mockDiscovery
	mockConfig := new(cli.MockCLIConfig)

	inputDir := "/some/input"
	outputDir := "/some/input/mp4"

	mockConfig.On("GetInputDir").Return(inputDir)
	mockConfig.On("GetOutputDir").Return(outputDir)

	mockFS.On("FileExists", inputDir).Return(true)
	mockDiscovery.On("CreateOutputDir", outputDir).Return(errors.New("disk error"))

	svc := services.NewHealthCheckService(mockDiscovery, mockFS, mockConfig)

	err := svc.Execute()

	assert.Error(t, err)
	mockConfig.AssertExpectations(t)
	mockFS.AssertExpectations(t)
	mockDiscovery.AssertExpectations(t)
}
