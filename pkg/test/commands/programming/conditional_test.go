package test

import (
	"testing"

	conditional "gorobot/pkg/commands/programming/conditional"
	console "gorobot/pkg/commands/programming/console"
	"gorobot/pkg/domain"
	"gorobot/pkg/engine"

	"github.com/stretchr/testify/assert"
)

func TestIfCommand_ShouldExecuteChildren_WhenConditionIsTrue(t *testing.T) {
	expressions := []string{
		"2 > 1",
		"true",
		`string.Equals("hello", "hello")`,
		"1+1 == 2",
	}

	for _, expr := range expressions {
		//Arange
		ifCmd := conditional.NewIfCommand(expr)
		write := console.NewWriteCommand("Hello, World!")

		ifCmd.AddCommand(write)
		s := domain.NewScript()
		s.AddCommand(ifCmd)

		e := engine.NewEngine()
		e.SetScript(s)

		called := map[string]any{}

		e.OnCommandStarted(func(cmd domain.Command) {
			called["commandStarted"] = cmd.GetID()
		})
		e.OnOutput(func(cmd domain.Command, value any) {
			called["output"] = value
		})

		//Act
		e.Run(true)

		//Assert
		assert.Equal(t, 2, len(called), "called should be length 2")
		assert.Equal(t, write.Message, called["output"], "output should be correct")
		assert.Equal(t, write.ID, called["commandStarted"], "commandStarted should be correct")

		e.Clear()
	}
}

func TestIfCommand_ShouldNotExecuteChildren_WhenConditionIsFalse(t *testing.T) {
	expressions := []string{
		"2 < 1",
		"false",
		`string.Equals("hello", "world")`,
		"1-1 == 2",
	}

	for _, expr := range expressions {
		//Arange
		ifCmd := conditional.NewIfCommand(expr)
		write := console.NewWriteCommand("Hello, World!")

		ifCmd.AddCommand(write)
		s := domain.NewScript()
		s.AddCommand(ifCmd)

		e := engine.NewEngine()
		e.SetScript(s)

		called := map[string]any{}

		e.OnCommandStarted(func(cmd domain.Command) {
			called["commandStarted"] = cmd.GetID()
		})

		//Act
		e.Run(true)

		//Assert
		assert.Equal(t, ifCmd.ID, called["commandStarted"], "commandStarted should be correct")
		assert.Equal(t, 1, len(called), "called should be length 1")
		e.Clear()
	}
}

func TestIfCommand_ShouldDisableNextElse_WhenConditionIsTrue(t *testing.T) {
	//Arange
	ifCmd := conditional.NewIfCommand("true")
	elseIfCmd := conditional.NewElseIfCommand("true")
	elseCmd := conditional.NewElseCommand()

	s := domain.NewScript()
	s.AddCommands([]domain.Command{
		ifCmd, elseIfCmd, elseCmd,
	})

	e := engine.NewEngine()
	e.SetScript(s)

	called := map[string]any{}

	e.OnCommandStarted(func(cmd domain.Command) {
		called["commandStarted"] = cmd.GetID()
	})

	//Act
	e.Run(true)

	//Assert
	assert.Equal(t, 1, len(called), "called should be length 1")
	assert.Equal(t, ifCmd.ID, called["commandStarted"], "commandStarted should be correct")

	e.Clear()
}

func TestElseIfCommand_ShouldExecuteChildren_WhenConditionIsTrue(t *testing.T) {
	expressions := []string{
		"2 > 1",
		"true",
		`string.Equals("hello", "hello")`,
		"1+1 == 2",
	}

	for _, expr := range expressions {
		//Arange
		ifCmd := conditional.NewIfCommand("false")
		elseIfCmd := conditional.NewElseIfCommand(expr)
		write := console.NewWriteCommand("Hello, World!")

		elseIfCmd.AddCommand(write)
		s := domain.NewScript()
		s.AddCommands([]domain.Command{
			ifCmd, elseIfCmd,
		})

		e := engine.NewEngine()
		e.SetScript(s)

		called := map[string]any{}

		e.OnCommandStarted(func(cmd domain.Command) {
			println(cmd.GetTag() + " " + cmd.GetID())
			called["commandStarted"] = cmd.GetID()
		})
		e.OnOutput(func(cmd domain.Command, value any) {
			called["output"] = value
		})

		//Act
		e.Run(true)

		//Assert
		assert.Equal(t, 2, len(called), "called should be length 2")
		assert.Equal(t, write.Message, called["output"], "output should be correct")
		assert.Equal(t, write.ID, called["commandStarted"], "commandStarted should be correct")

		e.Clear()
	}
}

func TestElseIfCommand_ShouldNotExecuteChildren_WhenConditionIsFalse(t *testing.T) {
	expressions := []string{
		"2 < 1",
		"false",
		`string.Equals("hello", "world")`,
		"1-1 == 2",
	}

	for _, expr := range expressions {
		//Arange
		ifCmd := conditional.NewIfCommand("false")
		elseIfCmd := conditional.NewElseIfCommand(expr)
		write := console.NewWriteCommand("Hello, World!")

		elseIfCmd.AddCommand(write)
		s := domain.NewScript()
		s.AddCommands([]domain.Command{
			ifCmd, elseIfCmd,
		})

		e := engine.NewEngine()
		e.SetScript(s)

		called := map[string]any{}

		e.OnCommandStarted(func(cmd domain.Command) {
			called["commandStarted"] = cmd.GetID()
		})

		//Act
		e.Run(true)

		//Assert
		assert.Equal(t, elseIfCmd.ID, called["commandStarted"], "commandStarted should be correct")
		assert.Equal(t, 1, len(called), "called should be length 1")
		e.Clear()
	}
}

func TestElseIfCommand_ShouldDisableNextElse_WhenConditionIsTrue(t *testing.T) {
	//Arange
	ifCmd := conditional.NewIfCommand("false")
	elseIfCmd := conditional.NewElseIfCommand("true")
	elseIf2Cmd := conditional.NewElseIfCommand("true")
	elseCmd := conditional.NewElseCommand()

	s := domain.NewScript()
	s.AddCommands([]domain.Command{
		ifCmd, elseIfCmd, elseIf2Cmd, elseCmd,
	})

	e := engine.NewEngine()
	e.SetScript(s)

	called := map[string]any{}

	e.OnCommandStarted(func(cmd domain.Command) {
		called["commandStarted"] = cmd.GetID()
	})

	//Act
	e.Run(true)

	//Assert
	assert.Equal(t, 1, len(called), "called should be length 1")
	assert.Equal(t, elseIfCmd.ID, called["commandStarted"], "commandStarted should be correct")

	e.Clear()
}

func TestElseCommand_ShouldExecute_WhenIfCommandConditionIsFalse(t *testing.T) {
	//Arange
	ifCmd := conditional.NewIfCommand("false")
	write := console.NewWriteCommand("Hello, World!")
	elseCmd := conditional.NewElseCommand()
	elseCmd.AddCommand(write)

	s := domain.NewScript()
	s.AddCommands([]domain.Command{
		ifCmd, elseCmd,
	})

	e := engine.NewEngine()
	e.SetScript(s)

	called := map[string]any{}

	e.OnCommandStarted(func(cmd domain.Command) {
		called["commandStarted"] = cmd.GetID()
	})

	e.OnOutput(func(cmd domain.Command, value any) {
		called["output"] = value
	})

	//Act
	e.Run(true)

	//Assert
	assert.Equal(t, 2, len(called), "called should be length 2")
	assert.Equal(t, write.ID, called["commandStarted"], "commandStarted should be correct")
	assert.Equal(t, write.Message, called["output"], "output should be correct")
	e.Clear()
}

func TestElseCommand_ShouldExecute_WhenElseIfCommandConditionIsFalse(t *testing.T) {
	//Arange
	ifCmd := conditional.NewIfCommand("false")
	elseIfCmd := conditional.NewElseIfCommand("false")
	write := console.NewWriteCommand("Hello, World!")
	elseCmd := conditional.NewElseCommand()
	elseCmd.AddCommand(write)

	s := domain.NewScript()
	s.AddCommands([]domain.Command{
		ifCmd, elseIfCmd, elseCmd,
	})

	e := engine.NewEngine()
	e.SetScript(s)

	called := map[string]any{}

	e.OnCommandStarted(func(cmd domain.Command) {
		called["commandStarted"] = cmd.GetID()
	})

	e.OnOutput(func(cmd domain.Command, value any) {
		called["output"] = value
	})

	//Act
	e.Run(true)

	//Assert
	assert.Equal(t, 2, len(called), "called should be length 2")
	assert.Equal(t, write.ID, called["commandStarted"], "commandStarted should be correct")
	assert.Equal(t, write.Message, called["output"], "output should be correct")
	e.Clear()
}
