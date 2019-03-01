package main

import (
	"context"
	"time"
)

type StateContext struct {
	ctx               context.Context
	currentState      WatchdogState
	latestStateChange time.Time

	OnIdle        func(current WatchdogState) WatchdogState
	AfterCall     func(operation string, err error, duration time.Duration)
	BeforeCall    func(operation string, state WatchdogState) (WatchdogState, error)
	OnStateChange func(current WatchdogState, new WatchdogState, duration time.Duration)
}

func (stateCtx *StateContext) State() WatchdogState {
	return stateCtx.currentState
}

func (stateCtx *StateContext) Call(operation string, f func(ctx context.Context, ws WatchdogState) (WatchdogState, error)) (WatchdogState, error) {

	// Before call logic
	if stateCtx.BeforeCall != nil {

		ws, err := stateCtx.BeforeCall(operation, stateCtx.currentState)

		if err != nil {
			// Change our state right before
			stateCtx.currentState = ws
			return ws, err
		}
	}

	// Actual call
	beginTime := time.Now()
	state, err := f(stateCtx.ctx, stateCtx.currentState)
	executionTime := time.Since(beginTime)

	dateSince := 0 * time.Second

	// Calculate time since last state
	if stateCtx.latestStateChange.IsZero() && stateCtx.currentState != state {
		stateCtx.latestStateChange = time.Now()
		dateSince = 0 * time.Second
	} else if stateCtx.currentState != state {
		dateSince = time.Since(stateCtx.latestStateChange)
		stateCtx.latestStateChange = time.Now()
	}

	// After call
	if stateCtx.AfterCall != nil {
		stateCtx.AfterCall(operation, err, executionTime)
	}

	// State change managment
	if stateCtx.OnStateChange != nil && stateCtx.currentState != state {
		stateCtx.OnStateChange(stateCtx.currentState, state, dateSince)
	} else if stateCtx.OnIdle != nil && stateCtx.currentState == state && err == nil {
		return stateCtx.OnIdle(state), nil
	}

	// Change actual status
	stateCtx.currentState = state

	return state, err
}
