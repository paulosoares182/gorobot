package test

import (
	"fmt"
	console "gorobot/pkg/commands/programming/console"
	variable "gorobot/pkg/commands/variable/modify"
	"gorobot/pkg/domain"
	"gorobot/pkg/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

type varTest struct {
	vtype string
	name  string
	value string
}

func TestSetVariable_ShouldUpsertVariable_WhenCommandIsValid(t *testing.T) {
	vars := []varTest{
		{vtype: "string", name: "var", value: "Hello World"},
		{vtype: "int", name: "var", value: "123"},
		{vtype: "float32", name: "var", value: "3.14"},
		{vtype: "float64", name: "var", value: "3.14159"},
		{vtype: "bool", name: "var", value: "true"},
		{vtype: "path", name: "var", value: "/usr/local/bin"},
	}

	for _, v := range vars {
		// Arrange
		setCmd := variable.NewSetVariableCommand(v.vtype, v.name, v.value)
		writeCmd := console.NewWriteCommand("{var}")

		s := domain.NewScript()
		s.AddCommands([]domain.Command{
			setCmd, writeCmd,
		})

		e := engine.NewEngine()
		e.SetScript(s)
		e.SetVariable(v.name, "")

		called := map[string]any{}
		e.OnOutput(func(cmd domain.Command, value any) {
			called["output"] = value
		})

		// Act
		e.Run(false)

		// Assert
		got := e.ListVariable()[0]
		assert.Equal(t, setCmd.Name, got.Name)
		assert.Equal(t, setCmd.Value, fmt.Sprint(got.Value))
		assert.Equal(t, setCmd.Value, called["output"], "output should be called")
	}
}

func TestSetVariable_ShouldNotUpsertVariable_WhenCommandIsValid(t *testing.T) {
	vars := []varTest{
		{vtype: "xstring", name: "var", value: "Hello World"},
		{vtype: "int", name: "var@", value: "123"},
		{vtype: "float32", name: "_ var", value: "3.14"},
		{vtype: "float64", name: "{var}", value: "3.14159"},
		{vtype: "bool", name: "@var", value: "true"},
		{vtype: "path", name: "va.r", value: "/usr/local/bin"},
	}

	for _, v := range vars {
		// Arrange
		setCmd := variable.NewSetVariableCommand(v.vtype, v.name, v.value)
		writeCmd := console.NewWriteCommand("{var}")

		s := domain.NewScript()
		s.AddCommands([]domain.Command{
			setCmd, writeCmd,
		})

		e := engine.NewEngine()
		e.SetScript(s)

		called := map[string]any{}
		e.OnOutput(func(cmd domain.Command, value any) {
			called["output"] = value
		})

		// Act
		ok, err := e.Run(false)

		// Assert
		got := len(e.ListVariable())

		assert.Error(t, err)
		assert.False(t, ok)
		assert.Equal(t, 0, got, "list variables should be empty")
		assert.Equal(t, 0, len(called), "output should not be called")
	}
}

func TestDeleteVariable_ShouldRemoveVariable(t *testing.T) {
	// Arrange
	setCmd := variable.NewSetVariableCommand("string", "var", "Hello World!")
	writeCmd := console.NewWriteCommand("{var}")
	delCmd := variable.NewDeleteVariableCommand("var")
	writeCmd2 := console.NewWriteCommand("{var}")

	s := domain.NewScript()
	s.AddCommands([]domain.Command{
		setCmd, writeCmd, delCmd, writeCmd2,
	})

	e := engine.NewEngine()
	e.SetScript(s)

	called := map[string]any{}
	e.OnOutput(func(cmd domain.Command, value any) {
		called[cmd.GetID()] = value
	})

	// Act
	e.Run(false)

	// Assert
	got := len(e.ListVariable())
	assert.Equal(t, 0, got, "list variables should be empty")
	assert.Equal(t, 2, len(called))
	assert.Equal(t, setCmd.Value, called[writeCmd.GetID()], "output should be called with value")
	assert.Equal(t, writeCmd2.Message, called[writeCmd2.GetID()], "output should be called with value")
}
