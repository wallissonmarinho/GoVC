package ports

// CommandExecutorPort defines the contract for command execution
type CommandExecutorPort interface {
	Register(name string, cmd ServiceCommand)
	Execute(name string) error
}
