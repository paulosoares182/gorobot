package commands

import (
	"fmt"
	"gorobot/pkg/domain"
	"gorobot/pkg/domain/exceptions"
	"reflect"
	"strconv"
	"strings"
)

type SetVariableCommand struct {
	domain.ScriptCommand
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
	value any    `json:"-"`
}

var SetVariableCommandTag = reflect.TypeOf(SetVariableCommand{}).Name()

func NewSetVariableCommand(vType string, name string, value string) *SetVariableCommand {
	return &SetVariableCommand{
		ScriptCommand: DefaultSetVariableCommand().ScriptCommand,
		Type:          vType,
		Name:          name,
		Value:         value,
		value:         value,
	}
}

func DefaultSetVariableCommand() *SetVariableCommand {
	return &SetVariableCommand{
		ScriptCommand: domain.NewCommand(SetVariableCommandTag, false),
	}
}

func (c *SetVariableCommand) Run(e domain.Engine) (any, error) {
	n := e.ExtractAsString(c.Name)
	c.value = e.ExtractAsString(c.Value)

	if err := c.Validate(); err != nil {
		return nil, err
	}

	err := e.SetVariable(n, c.value)

	return nil, err
}

func (c *SetVariableCommand) Validate() error {
	var errs []string

	switch strings.ToLower(c.Type) {
	case "string":
		c.value = c.Value

	case "int":
		if v, err := strconv.Atoi(c.Value); err != nil {
			errs = append(errs, fmt.Sprintf(`cannot convert "%s" to int`, c.Value))
		} else {
			c.value = v
		}

	case "float32":
		if v, err := strconv.ParseFloat(c.Value, 32); err != nil {
			errs = append(errs, fmt.Sprintf(`cannot convert "%s" to float32`, c.Value))
		} else {
			c.value = float32(v)
		}
	case "float64":
		if v, err := strconv.ParseFloat(c.Value, 64); err != nil {
			errs = append(errs, fmt.Sprintf(`cannot convert "%s" to float64`, c.Value))
		} else {
			c.value = float64(v)
		}
	case "bool":
		if v, err := strconv.ParseBool(c.Value); err != nil {
			errs = append(errs, fmt.Sprintf(`cannot convert "%s" to bool`, c.Value))
		} else {
			c.value = v
		}

	case "path":
		c.value = c.Value
	default:
		errs = append(errs, `field "Type" must be one of the following: string, int, float32, bool, path`)
	}

	if len(errs) > 0 {
		return exceptions.CommandError{Messages: errs}
	}

	return nil
}
