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

func (c *ScriptCommand) GetID() string               { return c.ID }
func (c *ScriptCommand) GetTag() string              { return c.Tag }
func (c *ScriptCommand) GetComment() string          { return c.Comment }
func (c *ScriptCommand) SetComment(comment string)   { c.Comment = comment }
func (c *ScriptCommand) CanHaveChildCommands() bool  { return c.CanHaveChildren }
func (c *ScriptCommand) GetCommands() []Command      { return c.Commands }
func (c *ScriptCommand) SetParentID(parentID string) { c.ParentID = parentID }

func (c *ScriptCommand) AddCommand(cmd Command) error {
	if !c.CanHaveChildren {
		return errors.New("this command cannot have child commands")
	}
	cmd.SetParentID(c.ID)
	c.Commands = append(c.Commands, cmd)
	return nil
}
