package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateArithmetic(t *testing.T) {
	out := Calculate("1+2")
	assert.Equal(t, "3", out)

	out = Calculate("10/2")
	assert.Equal(t, "5", out)

	out = Calculate("2+3*4")
	assert.Equal(t, "20", out)
}

func TestCalculateComparisonAndLogical(t *testing.T) {
	out := Calculate("5>3")
	assert.Equal(t, "true", out)

	out = Calculate("1<2 && 3>2")
	assert.Equal(t, "true", out)

	out = Calculate("(1+2)*3")
	assert.Equal(t, "9", out)
}
