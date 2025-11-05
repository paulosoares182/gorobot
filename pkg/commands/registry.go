package commands

import (
	actions "gorobot/pkg/commands/programming/actions"
	conditional "gorobot/pkg/commands/programming/conditional"
	console "gorobot/pkg/commands/programming/console"
	"gorobot/pkg/domain"
)

var CommandRegistry = map[string]func() domain.Command{}

func register(tag string, ctor func() domain.Command) {
	CommandRegistry[tag] = ctor
}

func init() {
	register(actions.CreateActionCommandTag, func() domain.Command { return actions.DefaultCreateActionCommand() })
	register(actions.InvokeActionCommandTag, func() domain.Command { return actions.DefaultInvokeActionCommand() })
	register(conditional.IfCommandTag, func() domain.Command { return conditional.DefaultIfCommand() })
	register(conditional.ElseIfCommandTag, func() domain.Command { return conditional.DefaultElseIfCommand() })
	register(conditional.ElseCommandTag, func() domain.Command { return conditional.DefaultElseCommand() })
	register(console.WriteCommandTag, func() domain.Command { return console.DefaultWriteCommand() })
}
