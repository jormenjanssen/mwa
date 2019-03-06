package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type WatchdogState int

const (
	Exit          WatchdogState = -1
	Preactivated  WatchdogState = 0
	Activated     WatchdogState = 1
	Alarm         WatchdogState = 2
	RestoreOk     WatchdogState = 3
	RestoreFailed WatchdogState = 4
)

func TranslateWatchdogState(state WatchdogState) string {

	switch state {

	case Exit:
		return "Exit"
	case Preactivated:
		return "Preactivated"
	case Activated:
		return "Activated"
	case Alarm:
		return "Alarm"
	case RestoreOk:
		return "Restored"
	case RestoreFailed:
		return "Not Recovered"
	default:
		return fmt.Sprintf("Not translated state: %v", state)
	}
}

func Watchdog(rctx RunContext, nhs NetworkHealthService) {

	log.Printf("Starting NetworkHealth Watchdog against target: %v with recovery time: [%v]", nhs.Address, nhs.RecoveryTime)

	// Start deactivated via function call, this helps with timing state changes and correct logging behaviour
	state, _ := ActivateWhenNoErrors(rctx, nhs)

	for {

		switch state {

		case Exit:
			return
		case Preactivated:
			state, _ = ActivateWhenNoErrors(rctx, nhs)
			break
		case Activated:
			state, _ = ActivateAlarmOnErrors(rctx, nhs)
			break
		case Alarm:
			state, _ = HandleAlarm(rctx, nhs, nhs)
			break
		case RestoreOk:
			state, _ = HandleRestored(rctx)
			break
		case RestoreFailed:
			state, _ = HandleRestoreFailure(rctx)
			break
		}
	}
}

func ActivateWhenNoErrors(rtcx RunContext, v Verify) (WatchdogState, error) {

	return rtcx.Call("verify", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := v.Verify(VerifyAttempts)

		if vErr != nil {
			return Preactivated, vErr
		}

		return Activated, nil
	})
}

func ActivateAlarmOnErrors(rtcx RunContext, v Verify) (WatchdogState, error) {

	return rtcx.Call("verify", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := v.Verify(VerifyAttempts)

		if vErr != nil {
			return Alarm, vErr
		}

		return Activated, nil
	})
}

func HandleAlarm(rtcx RunContext, r Recover, v Verify) (WatchdogState, error) {

	return rtcx.Call("recover", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		// Best effort restore
		r.Recover()

		// After recovery verify and increase the default number of attempts by 300%
		vErr := v.Verify(VerifyAttempts * 3)

		// Handle the restored case
		if vErr == nil {
			return RestoreOk, nil
		}

		return RestoreFailed, vErr
	})
}

func HandleRestored(rtcx RunContext) (WatchdogState, error) {

	return rtcx.Call("restored", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {
		return Activated, nil
	})
}

func HandleRestoreFailure(rtcx RunContext) (WatchdogState, error) {

	return rtcx.Call("restore-failed", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		// Wait a while before retrying
		<-time.After(25 * time.Second)

		return Activated, nil
	})
}
