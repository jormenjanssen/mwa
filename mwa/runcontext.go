package main

import "context"

type RunContext interface {
	State() WatchdogState
	Call(string, func(ctx context.Context, ws WatchdogState) (WatchdogState, error)) (WatchdogState, error)
}
