package utils

import (
	"strings"
)

func NormalizeLineBreaks(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func NormalizeStringLineBreaks(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, `\r\n`, "\n"), `\n`, "\n")
}

func RemoveNewLines(s string, replacement string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", replacement), "\n", replacement)
}

func RemoveTabs(s string, replacement string) string {
	return strings.ReplaceAll(s, "\t", replacement)
}

func RemoveVariableSyntax(text string) string {
	t := strings.TrimSpace(text)
	if strings.HasPrefix(t, "{") && strings.HasSuffix(t, "}") {
		return strings.TrimSuffix(strings.TrimPrefix(t, "{"), "}")
	}
	return text
}

func RemoveExpressionSyntax(text string) string {
	t := strings.TrimSpace(text)
	if strings.HasPrefix(t, "${") && strings.HasSuffix(t, "}") {
		return strings.TrimSuffix(strings.TrimPrefix(t, "${"), "}")
	}
	return text
}

func Substring(s string, maxLength int, suffix ...string) string {
	_suffix := "..."
	if len(suffix) > 0 {
		_suffix = suffix[0]
	}
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-len(_suffix)] + _suffix
}
