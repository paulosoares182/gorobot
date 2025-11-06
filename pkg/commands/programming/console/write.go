package commands

import (
	"gorobot/pkg/domain"
	"reflect"
)

type WriteCommand struct {
	domain.ScriptCommand
	Message string `json:"value"`
}

var WriteCommandTag = reflect.TypeOf(WriteCommand{}).Name()

func NewWriteCommand(message string) *WriteCommand {
	return &WriteCommand{
		ScriptCommand: DefaultWriteCommand().ScriptCommand,
		Message:       message,
	}
}

func DefaultWriteCommand() *WriteCommand {
	return &WriteCommand{
		ScriptCommand: domain.NewCommand(WriteCommandTag, false),
	}
}

func (c *WriteCommand) Run(e domain.Engine) (any, error) {
	m := e.ExtractAsString(c.Message)
	return m, nil
}
