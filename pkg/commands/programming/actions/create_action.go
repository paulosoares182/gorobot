package commands

import (
	"gorobot/pkg/domain"
	"reflect"
)

type CreateActionCommand struct {
	domain.ScriptCommand
	Name       string  `json:"name"`
	Parameters *string `json:"parameters"`
}

var CreateActionCommandTag = reflect.TypeOf(CreateActionCommand{}).Name()

func NewCreateActionCommand(name string, parameters *string) *CreateActionCommand {
	return &CreateActionCommand{
		ScriptCommand: DefaultCreateActionCommand().ScriptCommand,
		Name:          name,
		Parameters:    parameters,
	}
}

func DefaultCreateActionCommand() *CreateActionCommand {
	return &CreateActionCommand{
		ScriptCommand: domain.NewCommand(CreateActionCommandTag, true),
	}
}

func (c *CreateActionCommand) Run(e domain.Engine) (any, error) {
	n := e.ExtractAsString(c.Name)

	action := func() {
		if len(c.Commands) > 0 {
			s := e.GetScript()
			s.EnableAllCommands(c.Commands)
			for _, child := range c.Commands {
				err := e.ExecuteCommand(child)
				if err != nil {
					break
				}
			}
		}
	}

	//TODO - handle parameters
	t := domain.NewActionTemplate(action, []domain.ActionArgs{})
	err := e.SetVariable(n, t)

	return nil, err
}
