package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLineBreaks_ShouldConvertCRLFtoLF(t *testing.T) {
	// Arrange
	input := "line1\r\nline2"

	// Act
	result := NormalizeLineBreaks(input)

	// Assert
	assert.Equal(t, "line1\nline2", result)
}

func TestNormalizeLineBreaks_ShouldKeepTextWithoutLineBreaksUnchanged(t *testing.T) {
	// Arrange
	input := "no breaks"

	// Act
	result := NormalizeLineBreaks(input)

	// Assert
	assert.Equal(t, input, result)
}

func TestNormalizeStringLineBreaks_ShouldConvertEscapedNewlines(t *testing.T) {
	// Arrange
	input1 := `line1\nline2`
	input2 := `line1\r\nline2`

	// Act
	result1 := NormalizeStringLineBreaks(input1)
	result2 := NormalizeStringLineBreaks(input2)

	// Assert
	assert.Equal(t, "line1\nline2", result1)
	assert.Equal(t, "line1\nline2", result2)
}

func TestNormalizeStringLineBreaks_ShouldNotChangePlainText(t *testing.T) {
	// Arrange
	input := "plain text"

	// Act
	result := NormalizeStringLineBreaks(input)

	// Assert
	assert.Equal(t, input, result)
}

func TestRemoveNewLines_ShouldReplaceNewLinesWithReplacement(t *testing.T) {
	// Arrange
	input1 := "line1\nline2"
	input2 := "line1\rline2"

	// Act
	result1 := RemoveNewLines(input1, "-")
	result2 := RemoveNewLines(input2, "-")

	// Assert
	assert.Equal(t, "line1-line2", result1)
	assert.Equal(t, input2, result2)
}

func TestRemoveNewLines_ShouldNotChangeTextWithoutNewlines(t *testing.T) {
	// Arrange
	input := "single line text"

	// Act
	result := RemoveNewLines(input, " ")

	// Assert
	assert.Equal(t, input, result)
}

func TestRemoveTabs_ShouldReplaceTabsWithReplacement(t *testing.T) {
	// Arrange
	input := "one\ttwo"

	// Act
	result := RemoveTabs(input, " ")

	// Assert
	assert.Equal(t, "one two", result)
}

func TestRemoveTabs_ShouldNotChangeTextWithoutTabs(t *testing.T) {
	// Arrange
	input := "no tabs here"

	// Act
	result := RemoveTabs(input, " ")

	// Assert
	assert.Equal(t, input, result)
}

func TestRemoveVariableSyntax_ShouldRemoveOuterBrackets(t *testing.T) {
	// Arrange
	input1 := "{variable}"
	input2 := "{variable_2}"

	// Act
	result1 := RemoveVariableSyntax(input1)
	result2 := RemoveVariableSyntax(input2)

	// Assert
	assert.Equal(t, "variable", result1)
	assert.Equal(t, "variable_2", result2)
}

func TestRemoveVariableSyntax_ShouldNotChangeIfMissingOuterBracket(t *testing.T) {
	// Arrange
	input1 := "variable} "
	input2 := " { variable"
	input3 := "variable"

	// Act
	result1 := RemoveVariableSyntax(input1)
	result2 := RemoveVariableSyntax(input2)
	result3 := RemoveVariableSyntax(input3)

	// Assert
	assert.Equal(t, input1, result1)
	assert.Equal(t, input2, result2)
	assert.Equal(t, input3, result3)
}

func TestRemoveExpressionSyntax_ShouldRemoveExpressionSyntax(t *testing.T) {
	// Arrange
	input1 := "${variable}"
	input2 := `${
	variable_2
	}`

	// Act
	result1 := RemoveExpressionSyntax(input1)
	result2 := RemoveExpressionSyntax(input2)

	// Assert
	assert.Equal(t, "variable", result1)
	assert.Equal(t, `
	variable_2
	`, result2)
}

func TestRemoveExpressionSyntax_ShouldNotChangeIfMissingSyntax(t *testing.T) {
	// Arrange
	input1 := "{variable}"
	input2 := "${ variable"
	input3 := "$variable}"

	// Act
	result1 := RemoveExpressionSyntax(input1)
	result2 := RemoveExpressionSyntax(input2)
	result3 := RemoveExpressionSyntax(input3)

	// Assert
	assert.Equal(t, input1, result1)
	assert.Equal(t, input2, result2)
	assert.Equal(t, input3, result3)
}

func TestSubstring_ShouldReturnOriginalIfShorterOrEqualThanLimit(t *testing.T) {
	// Arrange
	input1 := "Hello, World!"
	input2 := "Hello, World!"

	// Act
	result1 := Substring(input1, 14)
	result2 := Substring(input2, 13)

	// Assert
	assert.Equal(t, input1, result1)
	assert.Equal(t, input2, result2)
}

func TestSubstring_ShouldTruncateProperlyConsideringSuffixLength(t *testing.T) {
	// Arrange
	input := "Hello, World!"
	expected := "Hello!!!"

	// Act
	result := Substring(input, 8, "!!!")

	// Assert
	assert.Equal(t, expected, result)
}
