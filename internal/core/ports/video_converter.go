package ports

import "github.com/wallissonmarinho/GoVC/internal/core/domain"

// VideoConverterPort is the output port for converting videos.
// (Right side of hexagon: core → external tools)
type VideoConverterPort interface {
	ConvertWithProgress(video *domain.Video, inputDir string, progressCallback func(float64)) error
	GetDuration(videoPath string) (float64, error)
	HasExternalSubtitles(video *domain.Video, inputDir string) bool
}
