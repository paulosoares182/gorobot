package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariable_Validate_Success(t *testing.T) {
	variable := Variable{
		Name:  "testVar",
		Value: "someValue",
	}
	err := variable.Validate()
	assert.NoError(t, err)
}

func TestVariable_Validate_MissingName(t *testing.T) {
	variable := Variable{
		Value: "someValue",
	}
	err := variable.Validate()
	assert.Error(t, err)
	assert.Equal(t, `field "Name" is required`, err.Error())
}

func TestExtractVariableValue_SimpleReplacement(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
	}
	text := "Hello, {name}!"
	want := "Hello, John Doe!"

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_NestedVariableExists(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
		{Name: "foo", Value: "bar"},
		{Name: "_foo", Value: "foo"},
	}
	text := "Hello, {name}! Your product: {{_foo}}."
	want := "Hello, John Doe! Your product: bar."

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_NestedVariableUnknown(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
		{Name: "_foo", Value: "foo"},
	}
	text := "Hello, {name}! Your product: {{_foo}}."
	want := "Hello, John Doe! Your product: {foo}."

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_UnknownVariable(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
	}
	text := "Hello, {name}! Code: {unknown}."
	want := "Hello, John Doe! Code: {unknown}."

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_MixedVariables(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
		{Name: "_foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
		{Name: "product", Value: "t-shirt"},
	}
	text := "Hello, {name}! Product: {{_foo}}, item: {product}, code: {code}."
	want := "Hello, John Doe! Product: bar, item: t-shirt, code: {code}."

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_NoVariables(t *testing.T) {
	vars := []Variable{}
	text := "Hello World"
	want := "Hello World"

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_UnderscoreAndNumbers(t *testing.T) {
	vars := []Variable{
		{Name: "_foo1", Value: "value1"},
		{Name: "bar_2", Value: "value2"},
	}
	text := "Testing {_foo1} and {bar_2}."
	want := "Testing value1 and value2."

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_MoreNested(t *testing.T) {
	vars := []Variable{
		{Name: "super_foo_bar", Value: "_foo"},
		{Name: "_foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
	}
	text := "Deep: {{{{super_foo_bar}}}} {_foo}:{foo}"
	want := "Deep: {bar} foo:bar"

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_MoreNestedUnknown(t *testing.T) {
	vars := []Variable{
		{Name: "_foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
	}
	text := "Deep: {{{{unknown}}}} {_foo}:{foo}"
	want := "Deep: {{{{unknown}}}} foo:bar"

	got := ExtractVariableValue(vars, text)
	assert.Equal(t, want, got)
}

func TestUpsertVariable_NewVariable(t *testing.T) {
	vars := []Variable{
		{Name: "var1", Value: "value1"},
	}
	updatedVars := UpsertVariable(vars, "var2", 182.00)
	assert.Len(t, updatedVars, 2)
	assert.Equal(t, 182.00, updatedVars[1].Value)
}

func TestUpsertVariable_UpdateExistingVariable(t *testing.T) {
	vars := []Variable{
		{Name: "var1", Value: "value1"},
	}
	updatedVars := UpsertVariable(vars, "var1", 1912)
	assert.Len(t, updatedVars, 1)
	assert.Equal(t, 1912, updatedVars[0].Value)
}

func TestRemoveVariable_RemoveExistingVariable(t *testing.T) {
	vars := []Variable{
		{Name: "var1", Value: "value1"},
		{Name: "var2", Value: "value2"},
	}

	updatedVars := RemoveVariable(vars, "var1")
	assert.Len(t, updatedVars, 1)
	assert.Equal(t, "var2", updatedVars[0].Name)
}

func TestRemoveVariable_DontRemoveVariableIfDoesNotExist(t *testing.T) {
	vars := []Variable{
		{Name: "var1", Value: "value1"},
	}

	updatedVars := RemoveVariable(vars, "var2")
	assert.Len(t, updatedVars, 1)
	assert.Equal(t, "var1", updatedVars[0].Name)
}

func TestIsValidVariableSyntax_ReturnTrueIfValid(t *testing.T) {
	validNames := []string{
		"{VARNAME}",
		"{var_name}",
		"{var123}",
		"{var_123}",
		"{_var123}",
		"{123var_}",
	}

	for _, name := range validNames {
		assert.True(t, IsValidVariableSyntax(name), "Expected valid for name: %s", name)
	}
}

func TestIsValidVariableSyntax_ReturnFalseIfInvalid(t *testing.T) {
	invalidNames := []string{
		"VARNAME}",
		"{var_name",
		"{ var123 }",
		"{var-123}",
		"{ var123}",
		"{123var }",
	}

	for _, name := range invalidNames {
		assert.False(t, IsValidVariableSyntax(name), "Expected invalid for name: %s", name)
	}
}
