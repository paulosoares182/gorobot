package commands

import (
	"gorobot/pkg/domain"

	"github.com/google/uuid"
)

type WriteCommand struct {
	domain.ScriptCommand
	Message string `json:"value"`
}

func NewWriteCommand(message ...string) *WriteCommand {
	_message := ""
	if len(message) > 0 {
		_message = message[0]
	}
	return &WriteCommand{
		ScriptCommand: domain.ScriptCommand{
			ID:              uuid.NewString(),
			Tag:             "WriteCommand",
			CanHaveChildren: false,
		},
		Message: _message,
	}
}

func (c *WriteCommand) Run(engine domain.Engine) (any, error) {
	println(c.Message)
	return c.Message, nil
}
