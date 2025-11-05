package commands

import (
	actions "gorobot/pkg/commands/programming/actions"
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
	register(console.WriteCommandTag, func() domain.Command { return console.DefaultWriteCommand() })
}
