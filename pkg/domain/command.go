package domain

import "errors"

type Command interface {
	GetID() string
	GetTag() string
	GetComment() string
	SetComment(comment string)
	CanHaveChildCommands() bool
	GetCommands() []Command
	AddCommand(cmd Command) error
	Run(engine Engine) (any, error)
	SetParentID(parentID string)
}

type ScriptCommand struct {
	ID              string    `json:"id"`
	Tag             string    `json:"tag"`
	Comment         string    `json:"comment,omitempty"`
	CanHaveChildren bool      `json:"-"`
	Commands        []Command `json:"commands,omitempty"`
	ParentID        string    `json:"parentId,omitempty"`
}

func (b *ScriptCommand) GetID() string               { return b.ID }
func (b *ScriptCommand) GetTag() string              { return b.Tag }
func (b *ScriptCommand) GetComment() string          { return b.Comment }
func (b *ScriptCommand) SetComment(comment string)   { b.Comment = comment }
func (b *ScriptCommand) CanHaveChildCommands() bool  { return b.CanHaveChildren }
func (b *ScriptCommand) GetCommands() []Command      { return b.Commands }
func (b *ScriptCommand) SetParentID(parentID string) { b.ParentID = parentID }

func (b *ScriptCommand) AddCommand(cmd Command) error {
	if !b.CanHaveChildren {
		return errors.New("this command cannot have child commands")
	}
	cmd.SetParentID(b.ID)
	b.Commands = append(b.Commands, cmd)
	return nil
}
