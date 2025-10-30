package commands

import (
	"gorobot/pkg/domain"

	"github.com/google/uuid"
)

type CreateActionCommand struct {
	domain.ScriptCommand
}

func NewCreateActionCommand() *CreateActionCommand {
	return &CreateActionCommand{
		ScriptCommand: domain.ScriptCommand{
			ID:              uuid.NewString(),
			Tag:             "CreateActionCommand",
			CanHaveChildren: true,
		},
	}
}

func (c *CreateActionCommand) Run(engine domain.Engine) (any, error) {
	if len(c.Commands) > 0 {
		for _, child := range c.Commands {
			engine.ExecuteCommand(child)
		}
	}
	return nil, nil
}
