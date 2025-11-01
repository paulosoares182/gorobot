package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type Variable struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func ExtractVariableValue(text string, vars []Variable) string {
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

func findVariable(vars []Variable, name string) *Variable {
	for i := range vars {
		if vars[i].Name == name {
			return &vars[i]
		}
	}
	return nil
}
