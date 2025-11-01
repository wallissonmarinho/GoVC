package cli

import (
	"fmt"
	"log"

	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

// CommandExecutor manages command execution
type CommandExecutor struct {
	commands map[string]ports.ServiceCommand
}

var _ ports.CommandExecutorPort = (*CommandExecutor)(nil)

// NewCommandExecutor creates a new CommandExecutor
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{
		commands: make(map[string]ports.ServiceCommand),
	}
}

// Register adds a command to the executor
func (ce *CommandExecutor) Register(name string, cmd ports.ServiceCommand) {
	ce.commands[name] = cmd
}

// Execute runs a command by name
func (ce *CommandExecutor) Execute(name string) error {
	cmd := ce.commands[name]
	if cmd == nil {
		err := fmt.Errorf("unknown command: %s", name)
		log.Printf("❌ %v", err)
		return err
	}

	log.Printf("▶️  Starting %s...\n", cmd.Name())
	if err := cmd.Execute(); err != nil {
		log.Printf("❌ %s failed: %v", cmd.Name(), err)
		return err
	}
	log.Printf("✅ %s completed!\n", cmd.Name())
	return nil
}
