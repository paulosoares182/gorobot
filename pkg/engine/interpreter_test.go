package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringHelpers(t *testing.T) {
	it := Interpreter{}

	res, err := it.Run(`string.Contains("hello","ell")`)
	assert.NoError(t, err)
	s, _ := res.(string)
	assert.Equal(t, "true", s)

	res, err = it.Run(`string.Replace("abba","b","x")`)
	assert.NoError(t, err)
	s, _ = res.(string)
	assert.Equal(t, "\"axxa\"", s)

	res, err = it.Run(`string.Equals("AbC","abc",IgnoreCase)`)
	assert.NoError(t, err)
	s, _ = res.(string)
	assert.Equal(t, "true", s)
}

func TestExecuteExpressionAndTestCondition(t *testing.T) {
	out, err := ExecuteExpression("${1>0}")
	assert.NoError(t, err)
	s, _ := out.(string)
	assert.Equal(t, "true", s)

	ok, err := TestCondition("${1>0}")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestGetDateTimeNow(t *testing.T) {
	dt, err := GetDateTime("${NOW}")
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), dt, time.Second*5)
}
