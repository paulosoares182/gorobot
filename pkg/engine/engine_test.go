package engine

import (
	"sync"
	"testing"

	actions "gorobot/pkg/commands/programming/actions"
	console "gorobot/pkg/commands/programming/console"
	"gorobot/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestEngineRunEvents(t *testing.T) {
	// Arrange
	ca := actions.NewCreateActionCommand("myAction", nil)
	w := console.NewWriteCommand("Hello, World!")
	ca.AddCommand(w)
	s := domain.NewScript()
	s.Commands = append(s.Commands, ca)

	eng := NewEngine()
	eng.SetScript(s)

	var mu sync.Mutex
	called := map[string]any{}

	eng.OnScriptStarted(func(script *domain.Script) {
		mu.Lock()
		called["scriptStarted"] = true
		mu.Unlock()
	})
	eng.OnCommandStarted(func(cmd domain.Command) {
		mu.Lock()
		called["commandStarted"] = true
		mu.Unlock()
	})
	eng.OnOutput(func(cmd domain.Command, value any) {
		mu.Lock()
		called["output"] = true
		called["outputValue"] = value
		mu.Unlock()
	})
	eng.OnCommandFinished(func(cmd domain.Command) {
		mu.Lock()
		called["commandFinished"] = true
		mu.Unlock()
	})

	// Act
	ok, err := eng.Run(false)

	// Assert
	assert.NoError(t, err)
	assert.True(t, ok)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, true, called["scriptStarted"], "scriptStarted should be called")
	assert.Equal(t, true, called["commandStarted"], "commandStarted should be called")
	assert.Equal(t, true, called["output"], "output should be called")
	assert.Equal(t, w.Message, called["outputValue"], "outputValue should be correct")
	assert.Equal(t, true, called["commandFinished"], "commandFinished should be called")
}
