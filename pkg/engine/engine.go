package engine

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gorobot/pkg/domain"
	"gorobot/pkg/serializer"
)

type EngineImpl struct {
	mu sync.RWMutex

	running   bool
	script    *domain.Script
	variables []domain.Variable

	services map[string]any

	httpClient *http.Client

	scriptStarted   []domain.ScriptStartedHandler
	commandStarted  []domain.CommandStartedHandler
	commandFinished []domain.CommandFinishedHandler
	commandExcept   []domain.CommandExceptionHandler
	scriptExcept    []domain.ScriptExceptionHandler
	outputHandlers  []domain.OutputHandler
	scriptFinished  []domain.ScriptFinishedHandler
}

func NewEngine() *EngineImpl {
	return &EngineImpl{
		services:   make(map[string]any),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (e *EngineImpl) OnScriptStarted(handler domain.ScriptStartedHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.scriptStarted = append(e.scriptStarted, handler)
}
func (e *EngineImpl) OnCommandStarted(handler domain.CommandStartedHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.commandStarted = append(e.commandStarted, handler)
}
func (e *EngineImpl) OnCommandFinished(handler domain.CommandFinishedHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.commandFinished = append(e.commandFinished, handler)
}
func (e *EngineImpl) OnCommandException(handler domain.CommandExceptionHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.commandExcept = append(e.commandExcept, handler)
}
func (e *EngineImpl) OnScriptException(handler domain.ScriptExceptionHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.scriptExcept = append(e.scriptExcept, handler)
}
func (e *EngineImpl) OnOutput(handler domain.OutputHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.outputHandlers = append(e.outputHandlers, handler)
}
func (e *EngineImpl) OnScriptFinished(handler domain.ScriptFinishedHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.scriptFinished = append(e.scriptFinished, handler)
}

func (e *EngineImpl) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

func (e *EngineImpl) ListVariable() []domain.Variable {
	e.mu.RLock()
	defer e.mu.RUnlock()
	vars := make([]domain.Variable, len(e.variables))
	copy(vars, e.variables)
	return vars
}

func (e *EngineImpl) RegisterService(name string, svc any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.services == nil {
		e.services = make(map[string]any)
	}
	e.services[name] = svc
}

func (e *EngineImpl) GetService(name string) (any, bool) {
	e.mu.RLock()
	v, ok := e.services[name]
	e.mu.RUnlock()
	return v, ok
}

func (e *EngineImpl) SetScriptFromJSON(jsonStr string) error {
	s, err := serializer.UnmarshalScript([]byte(jsonStr))
	if err != nil {
		return err
	}
	e.SetScript(s)
	return nil
}

func (e *EngineImpl) SetScript(script *domain.Script) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.script = script
}

func (e *EngineImpl) GetScript() *domain.Script {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.script
}

func (e *EngineImpl) GetHttpClient() *http.Client {
	if e.httpClient == nil {
		e.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return e.httpClient
}

func (e *EngineImpl) Run(throwException bool) (bool, error) {
	e.mu.Lock()
	if e.script == nil {
		e.mu.Unlock()
		return false, fmt.Errorf("script is required")
	}

	e.running = true
	script := e.script
	e.mu.Unlock()

	for _, h := range e.scriptStarted {
		h(script)
	}

	for _, c := range script.Commands {
		if !e.running {
			break
		}
		if err := e.executeCommand(c); err != nil {
			for _, h := range e.scriptExcept {
				h(err)
			}
			if throwException {
				e.mu.Lock()
				e.running = false
				e.mu.Unlock()
				return false, err
			}

			panic(err)
		}
	}

	e.mu.Lock()
	e.running = false
	e.mu.Unlock()

	for _, h := range e.scriptFinished {
		h(script.ID, "")
	}
	return true, nil
}

func (e *EngineImpl) ExecuteCommand(cmd domain.Command) {
	_ = e.executeCommand(cmd)
}

func (e *EngineImpl) executeCommand(cmd domain.Command) error {
	for _, h := range e.commandStarted {
		h(cmd)
	}

	if !e.running {
		return nil
	}

	res, err := cmd.Run(e)
	if err != nil {
		for _, h := range e.commandExcept {
			h(cmd, err)
		}
		return err
	}

	for _, h := range e.commandFinished {
		h(cmd)
	}

	if res != nil {
		for _, h := range e.outputHandlers {
			h(cmd, res)
		}
	}

	return nil
}

func (e *EngineImpl) TestCondition(expression string) bool {
	return expression == "true"
}

func (e *EngineImpl) GetDateTime(expression string) (time.Time, error) {
	if expression == "" {
		return time.Now(), nil
	}
	if t, err := time.Parse(time.RFC3339, expression); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("it was not possible to get the date: %s", expression)
}

func (e *EngineImpl) ExecuteExpression(expression string) (any, error) {
	return expression, nil
}

func (e *EngineImpl) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.script = nil
	e.variables = nil
	e.services = make(map[string]any)
	e.scriptStarted = nil
	e.commandStarted = nil
	e.commandFinished = nil
	e.commandExcept = nil
	e.scriptExcept = nil
	e.outputHandlers = nil
	e.scriptFinished = nil
}
