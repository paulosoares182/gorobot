package engine

import (
	"errors"
	"gorobot/pkg/utils"
	"regexp"
	"strings"
)

func TestCondition(expr string) (bool, error) {
	if strings.TrimSpace(expr) == "" {
		return false, nil
	}

	expr = utils.RemoveNewLines(expr, "")

	res, err := ExecuteExpression(expr)
	if err != nil {
		return false, err
	}

	switch v := res.(type) {
	case bool:
		return v, nil
	case string:
		s := strings.TrimSpace(strings.ToLower(v))
		if s == "true" {
			return true, nil
		}
		if s == "false" {
			return false, nil
		}
		return false, errors.New("expression did not evaluate to a boolean")
	default:
		return false, errors.New("expression did not evaluate to a boolean")
	}
}

func ExecuteExpression(expr string) (any, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return expr, nil
	}

	re := regexp.MustCompile(`\$\{([\W\w]+)\}`)
	m := re.FindStringSubmatch(expr)
	if len(m) > 1 {
		it := Interpreter{}
		return it.Run(m[1])
	}

	return expr, nil
}
