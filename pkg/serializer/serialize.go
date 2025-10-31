package serializer

import (
	"encoding/json"
	"fmt"

	actions "gorobot/pkg/commands/programming/actions"
	console "gorobot/pkg/commands/programming/console"
	"gorobot/pkg/domain"
)

var commandRegistry = map[string]func() domain.Command{}

func RegisterCommand(tag string, ctor func() domain.Command) {
	commandRegistry[tag] = ctor
}

func init() {
	RegisterCommand("CreateActionCommand", func() domain.Command { return actions.NewCreateActionCommand() })
	RegisterCommand("WriteCommand", func() domain.Command { return console.NewWriteCommand() })
}

func MarshalScript(s *domain.Script) ([]byte, error) {
	type alias domain.Script
	var flat []json.RawMessage

	var walk func(cmd domain.Command) error
	walk = func(cmd domain.Command) error {
		b, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}
		delete(m, "commands")
		if v, ok := m["parentId"]; ok {
			var pid string
			if err := json.Unmarshal(v, &pid); err != nil {
				return err
			}
			if pid == "" {
				delete(m, "parentId")
			}
		}
		rb, err := json.Marshal(m)
		if err != nil {
			return err
		}
		flat = append(flat, rb)

		for _, ch := range cmd.GetCommands() {
			if err := walk(ch); err != nil {
				return err
			}
		}
		return nil
	}

	for _, c := range s.Commands {
		if err := walk(c); err != nil {
			return nil, err
		}
	}

	return json.Marshal(&struct {
		*alias
		Commands []json.RawMessage `json:"commands"`
	}{
		alias:    (*alias)(s),
		Commands: flat,
	})
}

func UnmarshalScript(data []byte) (*domain.Script, error) {
	var aux struct {
		ID       string            `json:"id"`
		Schedule string            `json:"schedule"`
		Commands []json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}
	res := &domain.Script{ID: aux.ID, Schedule: aux.Schedule, Commands: []domain.Command{}}

	type item struct {
		cmd      domain.Command
		parentId string
	}
	all := make([]item, 0, len(aux.Commands))
	byID := map[string]domain.Command{}

	for _, cmdData := range aux.Commands {
		var temp map[string]json.RawMessage
		if err := json.Unmarshal(cmdData, &temp); err != nil {
			return nil, err
		}
		var tag string
		if err := json.Unmarshal(temp["tag"], &tag); err != nil {
			return nil, err
		}
		var parentId string
		if v, ok := temp["parentId"]; ok {
			if err := json.Unmarshal(v, &parentId); err != nil {
				return nil, err
			}
		}

		ctor, ok := commandRegistry[tag]
		if !ok {
			return nil, fmt.Errorf("tipo de comando desconhecido: %s", tag)
		}
		cmd := ctor()
		if err := json.Unmarshal(cmdData, cmd); err != nil {
			return nil, err
		}
		all = append(all, item{cmd: cmd, parentId: parentId})
		byID[cmd.GetID()] = cmd
	}

	for _, it := range all {
		if it.parentId == "" {
			res.Commands = append(res.Commands, it.cmd)
			continue
		}
		parent, ok := byID[it.parentId]
		if !ok {
			return nil, fmt.Errorf("parent not found: %s", it.parentId)
		}
		if err := parent.AddCommand(it.cmd); err != nil {
			return nil, err
		}
	}

	return res, nil
}
