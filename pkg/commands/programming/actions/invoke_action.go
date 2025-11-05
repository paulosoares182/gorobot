package commands

import (
	"errors"
	"gorobot/pkg/domain"
	"reflect"

	"github.com/google/uuid"
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
		ScriptCommand: domain.ScriptCommand{
			ID:              uuid.NewString(),
			Tag:             InvokeActionCommandTag,
			CanHaveChildren: false,
		},
	}
}

func (c *InvokeActionCommand) Run(e domain.Engine) (any, error) {
	v := e.ExtractAsAny(c.Name)

	t, ok := v.(*domain.ActionTemplate)
	if !ok {
		return nil, errors.New("action not found: " + c.Name)
	}

	t.Action()

	return nil, nil
}
