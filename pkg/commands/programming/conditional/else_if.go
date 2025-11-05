package commands

import (
	"fmt"
	"gorobot/pkg/domain"
	"reflect"

	"github.com/google/uuid"
)

type ElseIfCommand struct {
	domain.ScriptCommand
	Expression string `json:"expression"`
}

var ElseIfCommandTag = reflect.TypeOf(ElseIfCommand{}).Name()

func NewElseIfCommand(expression string) *ElseIfCommand {
	return &ElseIfCommand{
		ScriptCommand: DefaultElseIfCommand().ScriptCommand,
		Expression:    expression,
	}
}

func DefaultElseIfCommand() *ElseIfCommand {
	return &ElseIfCommand{
		ScriptCommand: domain.ScriptCommand{
			ID:              uuid.NewString(),
			Tag:             ElseIfCommandTag,
			CanHaveChildren: true,
		},
	}
}

func (c *ElseIfCommand) Run(e domain.Engine) (any, error) {
	ok, err := e.TestCondition(fmt.Sprintf("${%s}", c.Expression))

	if ok {
		c.disableNextElse(e)

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
	return nil, err
}

func (c *ElseIfCommand) disableNextElse(e domain.Engine) {
	s := e.GetScript()
	cmds := s.Commands

	if p := s.GetParent(c.ParentID); p != nil {
		cmds = p.GetCommands()
	}

	canDisable := false
	for _, cmd := range cmds {
		if cmd.GetID() == c.ID {
			canDisable = true
			continue
		}

		if canDisable {
			if cmd.GetTag() == ElseIfCommandTag {
				if elseIfCmd, ok := cmd.(*ElseIfCommand); ok {
					elseIfCmd.Enabled = false
					continue
				}
			}
			if cmd.GetTag() == ElseCommandTag {
				if elseCmd, ok := cmd.(*ElseCommand); ok {
					elseCmd.Enabled = false
					continue
				}
			}

			break
		}
	}
}
