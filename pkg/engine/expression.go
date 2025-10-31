package engine

import (
	"errors"
	"regexp"
	"strings"
)

func TestCondition(expression string) (bool, error) {
	if strings.TrimSpace(expression) == "" {
		return false, nil
	}

	//TODO - move this cleanup to utils
	expression = strings.ReplaceAll(expression, "\r\n", "")
	expression = strings.ReplaceAll(expression, "\n", "")

	res, err := ExecuteExpression(expression)
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

func ExecuteExpression(expression string) (any, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return expression, nil
	}

	re := regexp.MustCompile(`\$\{([\W\w]+)\}`)
	m := re.FindStringSubmatch(expression)
	if len(m) > 1 {
		it := Interpreter{}
		return it.Run(m[1])
	}

	return expression, nil
}
