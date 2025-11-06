package commands

import (
	"gorobot/pkg/domain"
	"reflect"
)

type DeleteVariableCommand struct {
	domain.ScriptCommand
	VariableName string `json:"variableName"`
}

var DeleteVariableCommandTag = reflect.TypeOf(DeleteVariableCommand{}).Name()

func NewDeleteVariableCommand(variableName string) *DeleteVariableCommand {
	return &DeleteVariableCommand{
		ScriptCommand: DefaultDeleteVariableCommand().ScriptCommand,
		VariableName:  variableName,
	}
}

func DefaultDeleteVariableCommand() *DeleteVariableCommand {
	return &DeleteVariableCommand{
		ScriptCommand: domain.NewCommand(DeleteVariableCommandTag, false),
	}
}

func (c *DeleteVariableCommand) Run(e domain.Engine) (any, error) {
	n := e.ExtractAsString(c.VariableName)

	e.DeleteVariable(n)

	return nil, nil
}
