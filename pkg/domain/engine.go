package domain

import (
	"net/http"
	"time"
)

type CommandStartedHandler func(cmd Command)
type CommandFinishedHandler func(cmd Command)
type ScriptStartedHandler func(script *Script)
type ScriptExceptionHandler func(err error)
type CommandExceptionHandler func(cmd Command, err error)
type OutputHandler func(cmd Command, obj any)
type ScriptFinishedHandler func(scriptID string, scriptName string)

type Engine interface {
	OnScriptStarted(handler ScriptStartedHandler)
	OnCommandStarted(handler CommandStartedHandler)
	OnCommandFinished(handler CommandFinishedHandler)
	OnCommandException(handler CommandExceptionHandler)
	OnScriptException(handler ScriptExceptionHandler)
	OnOutput(handler OutputHandler)
	OnScriptFinished(handler ScriptFinishedHandler)

	IsRunning() bool
	ListVariable() []Variable

	RegisterService(name string, svc any)
	GetService(name string) (any, bool)

	Run(throwException bool) (bool, error)
	SetScriptFromJSON(jsonStr string) error
	SetScript(script *Script)
	ExecuteCommand(cmd Command) bool
	GetScript() *Script
	GetHttpClient() *http.Client

	TestCondition(expression string) bool
	GetDateTime(expression string) (time.Time, error)
	ExecuteExpression(expression string) (any, error)

	Clear()
}
