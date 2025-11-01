package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVariableValue_SimpleReplacement(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
	}
	template := "Hello, {name}!"
	want := "Hello, John Doe!"

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_NestedVariableExists(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
		{Name: "foo", Value: "bar"},
		{Name: "_foo", Value: "foo"},
	}
	template := "Hello, {name}! Your product: {{_foo}}."
	want := "Hello, John Doe! Your product: bar."

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_NestedVariableUnknown(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
		{Name: "_foo", Value: "foo"},
	}
	template := "Hello, {name}! Your product: {{_foo}}."
	want := "Hello, John Doe! Your product: {foo}."

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_UnknownVariable(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
	}
	template := "Hello, {name}! Code: {unknown}."
	want := "Hello, John Doe! Code: {unknown}."

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_MixedVariables(t *testing.T) {
	vars := []Variable{
		{Name: "name", Value: "John Doe"},
		{Name: "_foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
		{Name: "product", Value: "t-shirt"},
	}
	template := "Hello, {name}! Product: {{_foo}}, item: {product}, code: {code}."
	want := "Hello, John Doe! Product: bar, item: t-shirt, code: {code}."

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_NoVariables(t *testing.T) {
	vars := []Variable{}
	template := "Hello World"
	want := "Hello World"

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_UnderscoreAndNumbers(t *testing.T) {
	vars := []Variable{
		{Name: "_foo1", Value: "value1"},
		{Name: "bar_2", Value: "value2"},
	}
	template := "Testing {_foo1} and {bar_2}."
	want := "Testing value1 and value2."

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_MoreNested(t *testing.T) {
	vars := []Variable{
		{Name: "super_foo_bar", Value: "_foo"},
		{Name: "_foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
	}
	template := "Deep: {{{{super_foo_bar}}}} {_foo}:{foo}"
	want := "Deep: {bar} foo:bar"

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}

func TestExtractVariableValue_MoreNestedUnknown(t *testing.T) {
	vars := []Variable{
		{Name: "_foo", Value: "foo"},
		{Name: "foo", Value: "bar"},
	}
	template := "Deep: {{{{unknown}}}} {_foo}:{foo}"
	want := "Deep: {{{{unknown}}}} foo:bar"

	got := ExtractVariableValue(template, vars)
	assert.Equal(t, want, got)
}
