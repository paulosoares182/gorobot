package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	actions "gorobot/pkg/commands/programming/actions"
	console "gorobot/pkg/commands/programming/console"

	"gorobot/pkg/domain"
)

func TestCommandStructure(t *testing.T) {
	// Arrange
	parent := actions.NewCreateActionCommand("parentAction", nil)
	child := actions.NewCreateActionCommand("childAction", nil)
	grand := console.NewWriteCommand("")
	parent2 := actions.NewCreateActionCommand("parent2Action", nil)
	child2 := console.NewWriteCommand("")

	// Act
	err1 := parent.AddCommand(child)
	err2 := child.AddCommand(grand)
	err3 := parent2.AddCommand(child2)

	// Assert: ParentID set correctly
	assert.NoError(t, err1)
	assert.Equal(t, parent.ID, child.ParentID)
	assert.NoError(t, err2)
	assert.Equal(t, child.ID, grand.ParentID)
	assert.NoError(t, err3)
	assert.Equal(t, parent2.ID, child2.ParentID)
}
func TestMarshalUnmarshalScript(t *testing.T) {
	// Arrange
	parent := actions.NewCreateActionCommand("parentAction", nil)
	child := actions.NewCreateActionCommand("childAction", nil)

	grand := console.NewWriteCommand("")

	parent2 := actions.NewCreateActionCommand("parent2Action", nil)
	child2 := console.NewWriteCommand("")

	parent.AddCommand(child)
	child.AddCommand(grand)
	parent2.AddCommand(child2)

	s := domain.NewScript()
	s.Commands = append(s.Commands, parent)
	s.Commands = append(s.Commands, parent2)

	// Act
	data, err := MarshalScript(s)
	assert.NoError(t, err)

	loaded, err := UnmarshalScript(data)
	assert.NoError(t, err)

	//Assert
	assert.Len(t, loaded.Commands, 2)

	root, ok := loaded.Commands[0].(*actions.CreateActionCommand)
	assert.True(t, ok)
	assert.Equal(t, parent.GetTag(), root.GetTag())

	assert.Len(t, root.Commands, 1)
	mid, ok := root.Commands[0].(*actions.CreateActionCommand)
	assert.True(t, ok)
	assert.Equal(t, child.ID, mid.GetID())
	assert.Equal(t, parent.ID, mid.ParentID)

	assert.Len(t, mid.Commands, 1)
	g, ok := mid.Commands[0].(*console.WriteCommand)
	assert.True(t, ok)
	assert.Equal(t, grand.ID, g.GetID())
	assert.Equal(t, child.ID, g.ParentID)

	root2, ok := loaded.Commands[1].(*actions.CreateActionCommand)
	assert.True(t, ok)
	assert.Equal(t, parent.GetTag(), root2.GetTag())

	assert.Len(t, root2.Commands, 1)
	c2, ok := root2.Commands[0].(*console.WriteCommand)
	assert.True(t, ok)
	assert.Equal(t, child2.ID, c2.GetID())
	assert.Equal(t, parent2.ID, c2.ParentID)
}
