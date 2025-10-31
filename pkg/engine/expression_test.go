package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteExpression_RawWhenNoPattern(t *testing.T) {
	out, err := ExecuteExpression("plain text")
	assert.NoError(t, err)
	s, ok := out.(string)
	assert.True(t, ok)
	assert.Equal(t, "plain text", s)
}

func TestExecuteExpression_EvaluatesExpression(t *testing.T) {
	out, err := ExecuteExpression("${1+2}")
	assert.NoError(t, err)
	s, ok := out.(string)
	assert.True(t, ok)
	assert.Equal(t, "3", s)
}

func TestTestCondition_TrueFalse(t *testing.T) {
	ok, err := TestCondition("${1>0}")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = TestCondition("${1>2}")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestTestCondition_NewlinesAndSpaces(t *testing.T) {
	ok, err := TestCondition("\n  ${1>0}  \n")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestTestCondition_InvalidBooleanError(t *testing.T) {
	_, err := TestCondition("${1+1}")
	assert.Error(t, err)
}
