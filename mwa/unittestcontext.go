package main

import (
	"context"
	"testing"
)

type UnitTestContext struct {
	stateCtx *StateContext
}

type ExpectedStateChange struct {
	From WatchdogState
	To   WatchdogState
}

func OnStateChange(c WatchdogState, n WatchdogState) {

}

func (utc UnitTestContext) State() WatchdogState {
	return utc.stateCtx.State()
}

func (utc UnitTestContext) Call(operation string, f func(ctx context.Context, ws WatchdogState) (WatchdogState, error)) (WatchdogState, error) {
	return utc.stateCtx.Call(operation, f)
}

func WrapUnitTestFunc(t *testing.T) {

}

func OnIdle(state WatchdogState) WatchdogState {
	return Exit
}

func NewUnitTestContext() UnitTestContext {
	ctx := UnitTestContext{stateCtx: &StateContext{OnIdle: OnIdle}}
	return ctx
}

func NewStateChangeUnitTestContext(t *testing.T, stateChanges []ExpectedStateChange) UnitTestContext {
	ctx := UnitTestContext{stateCtx: &StateContext{}}

	return ctx
}
