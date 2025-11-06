package commands

import (
	"errors"
	"gorobot/pkg/domain"
	"reflect"
)

type InvokeActionCommand struct {
	domain.ScriptCommand
	Name       string  `json:"name"`
	Parameters *string `json:"parameters"`
}

var InvokeActionCommandTag = reflect.TypeOf(InvokeActionCommand{}).Name()

func NewInvokeActionCommand(name string, parameters *string) *InvokeActionCommand {
	return &InvokeActionCommand{
		ScriptCommand: DefaultInvokeActionCommand().ScriptCommand,
		Name:          name,
		Parameters:    parameters,
	}
}

func DefaultInvokeActionCommand() *InvokeActionCommand {
	return &InvokeActionCommand{
		ScriptCommand: domain.NewCommand(InvokeActionCommandTag, false),
	}
}

func (c *InvokeActionCommand) Run(e domain.Engine) (any, error) {
	n := e.ExtractAsString(c.Name)
	v := e.ExtractAsAny(n)

	t, ok := v.(*domain.ActionTemplate)
	if !ok {
		return nil, errors.New("action not found: " + n)
	}

	t.Action()

	return nil, nil
}
