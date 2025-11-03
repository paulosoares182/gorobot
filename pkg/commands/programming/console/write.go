package commands

import (
	"gorobot/pkg/domain"

	"github.com/google/uuid"
)

type WriteCommand struct {
	domain.ScriptCommand
	Message string `json:"value"`
}

func NewWriteCommand(message string) *WriteCommand {
	return &WriteCommand{
		ScriptCommand: DefaultWriteCommand().ScriptCommand,
		Message:       message,
	}
}

func DefaultWriteCommand() *WriteCommand {
	return &WriteCommand{
		ScriptCommand: domain.ScriptCommand{
			ID:              uuid.NewString(),
			Tag:             "WriteCommand",
			CanHaveChildren: false,
		},
	}
}

func (c *WriteCommand) Run(engine domain.Engine) (any, error) {
	println(c.Message)
	return c.Message, nil
}
