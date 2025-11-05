package commands

import (
	"gorobot/pkg/domain"
	"reflect"

	"github.com/google/uuid"
)

type ElseCommand struct {
	domain.ScriptCommand
}

var ElseCommandTag = reflect.TypeOf(ElseCommand{}).Name()

func NewElseCommand() *ElseCommand {
	return &ElseCommand{
		ScriptCommand: DefaultElseCommand().ScriptCommand,
	}
}

func DefaultElseCommand() *ElseCommand {
	return &ElseCommand{
		ScriptCommand: domain.ScriptCommand{
			ID:              uuid.NewString(),
			Tag:             ElseCommandTag,
			CanHaveChildren: true,
		},
	}
}

func (c *ElseCommand) Run(e domain.Engine) (any, error) {
	s := e.GetScript()
	s.EnableAllCommands(c.Commands)
	for _, child := range c.Commands {
		err := e.ExecuteCommand(child)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
