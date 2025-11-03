package commands

import (
	"gorobot/pkg/domain"

	"github.com/google/uuid"
)

type CreateActionCommand struct {
	domain.ScriptCommand
	Name       string  `json:"name"`
	Parameters *string `json:"parameters"`
}

func NewCreateActionCommand(name string, parameters *string) *CreateActionCommand {
	return &CreateActionCommand{
		ScriptCommand: DefaultWriteCommand().ScriptCommand,
		Name:          name,
		Parameters:    parameters,
	}
}

func DefaultWriteCommand() *CreateActionCommand {
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
			ok := engine.ExecuteCommand(child)
			if !ok {
				break
			}
		}
	}
	return nil, nil
}
