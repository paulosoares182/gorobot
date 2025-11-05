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
			CanHaveChildren: true,
		},
	}
}

func (c *InvokeActionCommand) Run(engine domain.Engine) (any, error) {
	v := engine.ExtractAsAny(c.Name)

	t, ok := v.(*domain.ActionTemplate)
	if !ok {
		return nil, errors.New("action not found: " + c.Name)
	}

	t.Action()

	return nil, nil
}
