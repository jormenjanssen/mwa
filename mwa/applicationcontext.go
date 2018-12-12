package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type ApplicationContext struct {
	stateCtx *StateContext
	Delay    time.Duration
}

func (appCtx ApplicationContext) State() WatchdogState {
	return appCtx.stateCtx.currentState
}

func (appCtx ApplicationContext) Call(operation string, f func(ctx context.Context, ws WatchdogState) (WatchdogState, error)) (WatchdogState, error) {
	ws, err := appCtx.stateCtx.Call(operation, f)
	<-time.After(appCtx.Delay)
	return ws, err
}

func NewApplicationContext() ApplicationContext {
	return ApplicationContext{stateCtx: &StateContext{
		OnStateChange: LogStateChange,
		BeforeCall:    DelayWatchdog,
		AfterCall:     LogCall}, Delay: 1 * time.Second}
}

func LogStateChange(c WatchdogState, n WatchdogState) {
	if c != n {
		log.Println(fmt.Printf("Application state changes from: %v to: %v", TranslateWatchdogState(c), TranslateWatchdogState(n)))
	}
}

func LogCall(operation string, err error, duration time.Duration) {
	if err != nil {
		Debug(func() {
			log.Println(fmt.Printf("Call failed: %v with error: %v took %v", operation, err, duration))
		})
	} else {
		Debug(func() {
			log.Println(fmt.Printf("Succesfully ran call: %v took %v", operation, duration))
		})
	}
}

func DelayWatchdog(operation string, state WatchdogState) (WatchdogState, error) {

	if state == Activated && operation == "verify" {
		Debug(func() {
			log.Println("Delaying operation verify")
		})
		<-time.After(15 * time.Second)
	}

	return state, nil
}
