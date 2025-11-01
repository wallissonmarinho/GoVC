package ports

// ServiceCommand defines the contract for executable services
type ServiceCommand interface {
	Execute() error
	Name() string
}
