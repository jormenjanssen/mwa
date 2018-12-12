package main

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.Infof("Watchdog state changes from: [%v] to: [%v]", TranslateWatchdogState(c), TranslateWatchdogState(n))
	}
}

func LogCall(operation string, err error, duration time.Duration) {
	if err != nil {
		log.Debugf("CALL: [%v] [%v] [%v] [FAIL]", strings.ToUpper(operation), duration, err)
	} else {
		log.Debugf("CALL: [%v] [%v] [OK]", strings.ToUpper(operation), duration)
	}
}

func DelayWatchdog(operation string, state WatchdogState) (WatchdogState, error) {

	if state == Activated && operation == "verify" {
		// Only used in development logging log.Debugf("Delaying operation verify")
		<-time.After(17 * time.Second)
	}

	return state, nil
}
