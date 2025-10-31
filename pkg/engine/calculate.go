package engine

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reAndOrKeyWords = regexp.MustCompile(`(&&|\|\|)`)
	reValidParens   = regexp.MustCompile(`\(([^()]+)\)`)
	reIsNumber      = regexp.MustCompile(`^([\d]+)$`)
	reIsBoolean     = regexp.MustCompile(`^(?i:true|false)$`)
)

var keywords = map[string]bool{
	"+": true, "-": true, "*": true, "/": true,
	">": true, "<": true, ">=": true, "<=": true,
	"==": true, "!=": true, "&&": true, "||": true,
	"true": true, "false": true,
}

func Calculate(expression string) string {
	expression = strings.TrimSpace(expression)

	expression = ResolveParenthesis(expression)

	var result []string

	refresh := func() {
		if len(result) == 0 {
			return
		}
		expression = strings.Join(result, "")
		result = result[:0]
	}

	parts := tokenize(expression)
	calculateMathOperations(parts, &result)
	refresh()

	parts = tokenize(expression)
	calculateEqualityOperators(parts, &result)
	refresh()

	for reAndOrKeyWords.MatchString(expression) {
		parts = tokenize(expression)
		calculateLogicalOperators(expression, parts, &result)
		refresh()
	}

	return expression
}

func filterEmpty(arr []string) []string {
	out := make([]string, 0, len(arr))
	for _, s := range arr {
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func tokenize(expr string) []string {
	re := regexp.MustCompile(`"(?:\\.|[^"])*"|\d+(?:\.\d+)?|>=|<=|==|!=|&&|\|\||[+\-*/<>()]|[A-Za-z_][A-Za-z0-9_]*`)
	raw := re.FindAllString(expr, -1)
	return filterEmpty(raw)
}

func calculateMathOperations(parts []string, result *[]string) {
	for i := 0; i < len(parts); i++ {
		word := strings.ToLower(strings.TrimSpace(parts[i]))

		if keywords[word] {
			switch word {
			case "+", "-", "*", "/":
				prev := getPreviousDouble(parts, i, result)
				next, _ := strconv.ParseFloat(strings.TrimSpace(parts[i+1]), 64)

				var x float64
				switch word {
				case "+":
					x = prev + next
				case "-":
					x = prev - next
				case "*":
					x = prev * next
				case "/":
					x = prev / next
				}

				s := strconv.FormatFloat(x, 'G', -1, 64)
				*result = append(*result, s)
				i++
			case ">", "<", ">=", "<=", "==", "!=", "&&", "||", "true", "false":
				*result = append(*result, word)
			}
		} else {
			if reIsNumber.MatchString(strings.TrimSpace(word)) {
				*result = append(*result, word)
			}
		}
	}
}

func calculateEqualityOperators(parts []string, result *[]string) {
	for i := 0; i < len(parts); i++ {
		word := strings.ToLower(strings.TrimSpace(parts[i]))

		if keywords[word] {
			switch word {
			case ">", "<", ">=", "<=", "==", "!=":
				x := compare(word, parts[i-1], parts[i+1], parts, i, result)
				*result = append(*result, strconv.FormatBool(x))
				i++
			case "&&", "||", "true", "false":
				*result = append(*result, word)
			}
		}
	}
}

func calculateLogicalOperators(expression string, parts []string, result *[]string) {
	for i := 0; i < len(parts); i++ {
		word := strings.ToLower(strings.TrimSpace(parts[i]))

		if keywords[word] {
			switch word {
			case "&&":
				prev := getPreviousBoolean(parts, i, result)
				next, _ := strconv.ParseBool(strings.TrimSpace(parts[i+1]))
				*result = append(*result, strconv.FormatBool(prev && next))
				i++
			case "||":
				prev := getPreviousBoolean(parts, i, result)
				next, _ := strconv.ParseBool(strings.TrimSpace(parts[i+1]))
				*result = append(*result, strconv.FormatBool(prev || next))
				i++
			}
		}
	}
}

func ResolveParenthesis(expression string) string {
	tempExpression := expression

	for {
		matches := reValidParens.FindAllStringSubmatch(tempExpression, -1)
		if len(matches) == 0 {
			break
		}
		for _, m := range matches {
			inner := strings.TrimSpace(m[1])
			res := Calculate(inner)
			tempExpression = strings.Replace(tempExpression, m[0], res, 1)
		}
	}

	return tempExpression
}

func compare(operator string, item1 string, item2 string, parts []string, currentIndex int, result *[]string) bool {
	if operator == "==" {
		if reIsBoolean.MatchString(strings.TrimSpace(item1)) {
			return getPreviousBoolean(parts, currentIndex, result) == parseBool(item2)
		}
		return getPreviousDouble(parts, currentIndex, result) == parseFloat(item2)
	} else if operator == "!=" {
		if reIsBoolean.MatchString(strings.TrimSpace(item1)) {
			return getPreviousBoolean(parts, currentIndex, result) != parseBool(item2)
		}
		return getPreviousDouble(parts, currentIndex, result) != parseFloat(item2)
	}

	switch operator {
	case ">":
		return getPreviousDouble(parts, currentIndex, result) > parseFloat(item2)
	case "<":
		return getPreviousDouble(parts, currentIndex, result) < parseFloat(item2)
	case ">=":
		return getPreviousDouble(parts, currentIndex, result) >= parseFloat(item2)
	case "<=":
		return getPreviousDouble(parts, currentIndex, result) <= parseFloat(item2)
	}

	return false
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}

func parseBool(s string) bool {
	v, _ := strconv.ParseBool(strings.TrimSpace(s))
	return v
}

func getPreviousDouble(parts []string, currentIndex int, result *[]string) float64 {
	if len(*result) > 0 {
		last := (*result)[len(*result)-1]
		if v, err := strconv.ParseFloat(last, 64); err == nil {
			*result = (*result)[:len(*result)-1]
			return v
		}
	}
	v, _ := strconv.ParseFloat(strings.TrimSpace(parts[currentIndex-1]), 64)
	return v
}

func getPreviousBoolean(parts []string, currentIndex int, result *[]string) bool {
	if len(*result) > 0 {
		last := (*result)[len(*result)-1]
		if v, err := strconv.ParseBool(last); err == nil {
			*result = (*result)[:len(*result)-1]
			return v
		}
	}
	v, _ := strconv.ParseBool(strings.TrimSpace(parts[currentIndex-1]))
	return v
}
