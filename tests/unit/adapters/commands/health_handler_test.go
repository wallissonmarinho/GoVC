package commands

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	urfavecli "github.com/urfave/cli/v2"
	"github.com/wallissonmarinho/GoVC/internal/adapters/commands"
	"github.com/wallissonmarinho/GoVC/internal/core/ports"
)

type fakeExecutor struct {
	called bool
}

func (f *fakeExecutor) Execute() error { f.called = true; return nil }

func TestHealthCommandHandler_Execute_CallsExecutor(t *testing.T) {
	fake := &fakeExecutor{}

	var factory commands.ExecutorFactory = func(c *urfavecli.Context) (ports.Executor, error) {
		return fake, nil
	}

	handler := commands.NewHealthCommandHandler(factory)

	err := handler.Execute(nil)

	assert.NoError(t, err)
	assert.True(t, fake.called, "executor should have been called")
}

func TestHealthCommandHandler_Execute_FactoryError(t *testing.T) {
	var factory commands.ExecutorFactory = func(c *urfavecli.Context) (ports.Executor, error) {
		return nil, errors.New("factory error")
	}

	handler := commands.NewHealthCommandHandler(factory)
	err := handler.Execute(nil)

	assert.Error(t, err)
}
