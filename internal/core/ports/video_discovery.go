package ports

import "github.com/wallissonmarinho/GoVC/internal/core/domain"

// VideoDiscoveryPort is the input port for discovering videos.
// (Left side of hexagon: external world → core)
type VideoDiscoveryPort interface {
	FindVideos(dir string) ([]*domain.Video, error)
	CreateOutputDir(dir string) error
}
