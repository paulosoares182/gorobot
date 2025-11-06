package domain

import (
	"errors"
	"fmt"
	"gorobot/pkg/utils"
	"regexp"
	"strings"
)

type Variable struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func (v Variable) Validate() error {
	if v.Name == "" {
		return errors.New(`field "Name" is required`)
	}
	if v.Value == nil {
		return errors.New(`field "Value" is required`)
	}
	return nil
}

func ExtractAsString(vars []Variable, text string) string {
	groupRe := regexp.MustCompile(`^(\{+)([a-zA-Z0-9_]+)(\}+)$`)
	findRe := regexp.MustCompile(`\{+[a-zA-Z0-9_]+\}+`)

	result := text
	for i := 0; i < 100; i++ {
		changed := false

		newResult := findRe.ReplaceAllStringFunc(result, func(match string) string {
			parts := groupRe.FindStringSubmatch(match)
			if len(parts) != 4 {
				return match
			}
			opens := parts[1]
			variableName := parts[2]
			closes := parts[3]

			k := len(opens)
			if len(closes) < k {
				k = len(closes)
			}

			if k == 1 {
				v := findVariable(vars, variableName)
				if v != nil {
					if !isPrimitiveType(v.Value) {
						return match
					}
					changed = true
					return fmt.Sprint(v.Value)
				}
				return match
			}

			if k > 1 {
				v := findVariable(vars, variableName)
				if v == nil {
					return match
				}

				if !isPrimitiveType(v.Value) {
					return match
				}

				vValue := fmt.Sprint(v.Value)

				replacement := strings.Repeat("{", k-1) + vValue + strings.Repeat("}", k-1)
				changed = true
				return replacement
			}

			return match
		})

		if !changed || newResult == result {
			break
		}
		result = newResult
	}

	return result
}

func ExtractAsAny(vars []Variable, name string) any {

	v := findVariable(vars, utils.RemoveVariableSyntax(strings.TrimSpace(name)))
	if v != nil {
		return v.Value
	}
	return nil
}

func UpsertVariable(vars []Variable, name string, value any) ([]Variable, error) {
	if ok := IsValidVariableName(name); !ok {
		return vars, fmt.Errorf("variable name '%s' is not valid", name)
	}

	v := findVariable(vars, name)
	if v != nil {
		v.Value = value
		return vars, nil
	}

	vars = append(vars, Variable{
		Name:  name,
		Value: value,
	})
	return vars, nil
}

func RemoveVariable(vars []Variable, name string) []Variable {
	for i := range vars {
		if vars[i].Name == name {
			return append(vars[:i], vars[i+1:]...)
		}
	}
	return vars
}

func IsValidVariableSyntax(text string) bool {
	return regexp.MustCompile(`^\{[a-zA-Z0-9_]+\}$`).MatchString(strings.TrimSpace(text))
}

func IsValidVariableName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(strings.TrimSpace(name))
}

func isPrimitiveType(v any) bool {
	switch v.(type) {
	case string, int, float32, float64, bool:
		return true
	default:
		return false
	}
}

func findVariable(vars []Variable, name string) *Variable {
	for i := range vars {
		if vars[i].Name == name {
			return &vars[i]
		}
	}
	return nil
}
