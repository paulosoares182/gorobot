package domain

type Script struct {
	ID       string    `json:"id"`
	Schedule string    `json:"schedule,omitempty"`
	Commands []Command `json:"commands"`
}

func NewScript() *Script {
	return &Script{
		ID:       "",
		Commands: []Command{},
	}
}

func (s *Script) AddCommand(cmd Command) error {
	s.Commands = append(s.Commands, cmd)
	return nil
}

func (s *Script) AddCommands(cmds []Command) error {
	for _, cmd := range cmds {
		s.AddCommand(cmd)
	}
	return nil
}

func (s *Script) GetParent(parentID string) Command {
	if parentID == "" {
		return nil
	}
	return findCommand(s.Commands, parentID)
}

func (s *Script) EnableAllCommands(cmds []Command, recursive ...bool) {
	r := false
	if len(recursive) > 0 {
		r = recursive[0]
	}

	if cmds == nil {
		enableCommands(s.Commands, r)
	} else {
		enableCommands(cmds, r)
	}
}

func enableCommands(cmds []Command, recursive bool) {

	for _, cmd := range cmds {
		cmd.SetEnabled(true)
		if recursive && cmd.CanHaveChildCommands() && len(cmd.GetCommands()) > 0 {
			enableCommands(cmd.GetCommands(), recursive)
		}
	}
}

func findCommand(cmds []Command, parentID string) Command {
	for _, cmd := range cmds {
		if cmd.GetID() == parentID {
			return cmd
		}
		if cmd.CanHaveChildCommands() && len(cmd.GetCommands()) > 0 {
			if parent := findCommand(cmd.GetCommands(), parentID); parent != nil {
				return parent
			}
		}
	}
	return nil
}
