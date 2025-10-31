package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Interpreter struct{}

func (i *Interpreter) Run(expression string) (any, error) {
	expr := i.evaluateMethods(expression)

	//TODO - move this cleanup to utils
	expr = strings.ReplaceAll(expr, "\r\n", "")
	expr = strings.ReplaceAll(expr, "\n", "")

	res := Calculate(expr)

	return res, nil
}

func (i *Interpreter) evaluateMethods(expression string) string {
	expr := expression

	// ToUpper
	reToUpper := regexp.MustCompile(`string\.ToUpper\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reToUpper.ReplaceAllStringFunc(expr, func(s string) string {
		m := reToUpper.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return `"` + strings.ToUpper(m[1]) + `"`
	})

	// ToLower
	reToLower := regexp.MustCompile(`string\.ToLower\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reToLower.ReplaceAllStringFunc(expr, func(s string) string {
		m := reToLower.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return `"` + strings.ToLower(m[1]) + `"`
	})

	// Trim
	reTrim := regexp.MustCompile(`string\.Trim\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reTrim.ReplaceAllStringFunc(expr, func(s string) string {
		m := reTrim.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return `"` + strings.TrimSpace(m[1]) + `"`
	})

	// Length
	reLen := regexp.MustCompile(`string\.Length\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reLen.ReplaceAllStringFunc(expr, func(s string) string {
		m := reLen.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return strconv.Itoa(len(m[1]))
	})

	// IsNullOrWhiteSpace / IsNullOrEmpty and negations
	reIsNullOrWhiteSpace := regexp.MustCompile(`string\.IsNullOrWhiteSpace\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reIsNullOrWhiteSpace.ReplaceAllStringFunc(expr, func(s string) string {
		m := reIsNullOrWhiteSpace.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return strconv.FormatBool(strings.TrimSpace(m[1]) == "")
	})

	reIsNullOrEmpty := regexp.MustCompile(`string\.IsNullOrEmpty\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reIsNullOrEmpty.ReplaceAllStringFunc(expr, func(s string) string {
		m := reIsNullOrEmpty.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return strconv.FormatBool(m[1] == "")
	})

	reIsNotNullOrWhiteSpace := regexp.MustCompile(`string\.IsNotNullOrWhiteSpace\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reIsNotNullOrWhiteSpace.ReplaceAllStringFunc(expr, func(s string) string {
		m := reIsNotNullOrWhiteSpace.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return strconv.FormatBool(strings.TrimSpace(m[1]) != "")
	})

	reIsNotNullOrEmpty := regexp.MustCompile(`string\.IsNotNullOrEmpty\(\s*"((?:.|\\n)*?)"\s*\)`)
	expr = reIsNotNullOrEmpty.ReplaceAllStringFunc(expr, func(s string) string {
		m := reIsNotNullOrEmpty.FindStringSubmatch(s)
		if len(m) < 2 {
			return s
		}
		return strconv.FormatBool(m[1] != "")
	})

	// Contains / DoesNotContain: string.Contains("a","b",IgnoreCase?)
	reContains := regexp.MustCompile(`string\.Contains\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reContains.ReplaceAllStringFunc(expr, func(s string) string {
		m := reContains.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(strings.Contains(strings.ToLower(a), strings.ToLower(b)))
		}
		return strconv.FormatBool(strings.Contains(a, b))
	})

	reDoesNotContain := regexp.MustCompile(`string\.DoesNotContain\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reDoesNotContain.ReplaceAllStringFunc(expr, func(s string) string {
		m := reDoesNotContain.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(!strings.Contains(strings.ToLower(a), strings.ToLower(b)))
		}
		return strconv.FormatBool(!strings.Contains(a, b))
	})

	// Equals / IsNotEqual
	reEquals := regexp.MustCompile(`string\.Equals\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reEquals.ReplaceAllStringFunc(expr, func(s string) string {
		m := reEquals.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(strings.EqualFold(a, b))
		}
		return strconv.FormatBool(a == b)
	})

	reIsNotEqual := regexp.MustCompile(`string\.IsNotEqual\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reIsNotEqual.ReplaceAllStringFunc(expr, func(s string) string {
		m := reIsNotEqual.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(!strings.EqualFold(a, b))
		}
		return strconv.FormatBool(a != b)
	})

	// StartsWith / DoesNotStartWith
	reStartsWith := regexp.MustCompile(`string\.StartsWith\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reStartsWith.ReplaceAllStringFunc(expr, func(s string) string {
		m := reStartsWith.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(strings.HasPrefix(strings.ToLower(a), strings.ToLower(b)))
		}
		return strconv.FormatBool(strings.HasPrefix(a, b))
	})

	reDoesNotStartWith := regexp.MustCompile(`string\.DoesNotStartWith\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reDoesNotStartWith.ReplaceAllStringFunc(expr, func(s string) string {
		m := reDoesNotStartWith.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(!strings.HasPrefix(strings.ToLower(a), strings.ToLower(b)))
		}
		return strconv.FormatBool(!strings.HasPrefix(a, b))
	})

	// EndsWith / DoesNotEndWith
	reEndsWith := regexp.MustCompile(`string\.EndsWith\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reEndsWith.ReplaceAllStringFunc(expr, func(s string) string {
		m := reEndsWith.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(strings.HasSuffix(strings.ToLower(a), strings.ToLower(b)))
		}
		return strconv.FormatBool(strings.HasSuffix(a, b))
	})

	reDoesNotEndWith := regexp.MustCompile(`string\.DoesNotEndWith\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reDoesNotEndWith.ReplaceAllStringFunc(expr, func(s string) string {
		m := reDoesNotEndWith.FindStringSubmatch(s)
		if len(m) < 3 {
			return s
		}
		a := m[1]
		b := m[2]
		if len(m) >= 4 && (strings.EqualFold(m[3], "IgnoreCase") || strings.EqualFold(m[3], "IC")) {
			return strconv.FormatBool(!strings.HasSuffix(strings.ToLower(a), strings.ToLower(b)))
		}
		return strconv.FormatBool(!strings.HasSuffix(a, b))
	})

	// Replace: string.Replace("a","b","c",IgnoreCase?) -> returns quoted replaced string
	reReplace := regexp.MustCompile(`string\.Replace\(\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*,\s*"((?:.|\\n)*?)"\s*(?:,\s*(IgnoreCase|IC))?\s*\)`)
	expr = reReplace.ReplaceAllStringFunc(expr, func(s string) string {
		m := reReplace.FindStringSubmatch(s)
		if len(m) < 4 {
			return s
		}
		a := m[1]
		old := m[2]
		newv := m[3]
		if len(m) >= 5 && (strings.EqualFold(m[4], "IgnoreCase") || strings.EqualFold(m[4], "IC")) {
			rr := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(old))
			return `"` + rr.ReplaceAllString(a, newv) + `"`
		}
		return `"` + strings.ReplaceAll(a, old, newv) + `"`
	})

	return expr
}

func GetDateTime(expression string) (time.Time, error) {
	expression = strings.TrimSpace(expression)

	exprNoSpace := strings.ReplaceAll(expression, " ", "")

	if strings.EqualFold(exprNoSpace, "${NOW}") {
		return time.Now(), nil
	}
	if strings.EqualFold(exprNoSpace, "${UTC_NOW}") {
		return time.Now().UTC(), nil
	}

	if strings.Contains(strings.ToLower(expression), "parse") {
		return parseDateTime(expression)
	}

	parts := strings.Split(expression, ",")
	ints := []int{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid integer in datetime parts: %w", err)
		}
		ints = append(ints, v)
	}

	switch len(ints) {
	case 3:
		return time.Date(ints[0], time.Month(ints[1]), ints[2], 0, 0, 0, 0, time.Local), nil
	case 4:
		return time.Date(ints[0], time.Month(ints[1]), ints[2], ints[3], 0, 0, 0, time.Local), nil
	case 5:
		return time.Date(ints[0], time.Month(ints[1]), ints[2], ints[3], ints[4], 0, 0, time.Local), nil
	case 6:
		return time.Date(ints[0], time.Month(ints[1]), ints[2], ints[3], ints[4], ints[5], 0, time.Local), nil
	default:
		return time.Time{}, fmt.Errorf("%s is invalid", expression)
	}
}

func parseDateTime(expression string) (time.Time, error) {
	re := regexp.MustCompile(`"((?:.|\n)*?)"`)
	m := re.FindAllStringSubmatch(expression, -1)
	if len(m) < 2 {
		return time.Time{}, fmt.Errorf("%s is invalid. Minimum 2 arguments required", expression)
	}
	dateStr := m[0][1]
	format := m[1][1]

	layout := goFormat(format)

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s is invalid. %v", expression, err)
	}
	return t, nil
}

func goFormat(format string) string {
	replacements := []struct {
		from string
		to   string
	}{
		{"yyyy", "2006"},
		{"yyy", "2006"},
		{"yy", "06"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"MM", "01"},
		{"M", "1"},
		{"dddd", "Monday"},
		{"ddd", "Mon"},
		{"dd", "02"},
		{"d", "2"},
		{"HH", "15"},
		{"H", "15"},
		{"hh", "03"},
		{"h", "3"},
		{"mm", "04"},
		{"m", "4"},
		{"ss", "05"},
		{"s", "5"},
		{"fff", "000"},
		{"ff", "00"},
		{"f", "0"},
		{"tt", "PM"},
		{"t", "P"},
	}

	layout := format

	for _, r := range replacements {
		layout = strings.ReplaceAll(layout, r.from, r.to)
	}

	if layout == "" {
		return format
	}

	return layout
}
