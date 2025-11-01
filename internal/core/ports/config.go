package ports

// ConfigPort is the input port for configuration.
// (Left side of hexagon: external world → core)
type ConfigPort interface {
	GetInputDir() string
	GetOutputDir() string
	GetWorkers() int
	SaveLogsEnabled() bool
}
