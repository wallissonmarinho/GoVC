package ports

// Executor defines what a service can do (generic interface)
type Executor interface {
	Execute() error
}
