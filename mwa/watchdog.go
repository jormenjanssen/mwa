package main

import (
	"context"
	"fmt"
)

type WatchdogState int

const (
	Exit         WatchdogState = -1
	Preactivated WatchdogState = 0
	Activated    WatchdogState = 1
	Alarm        WatchdogState = 2
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
	default:
		return fmt.Sprintf("Not translated state: %v", state)
	}
}

func Watchdog(rctx RunContext, v Verify, r Recover) {

	state := Preactivated

	for {

		switch state {

		case Exit:
			return
		case Preactivated:
			state, _ = ActivateWhenNoErrors(rctx, v)
			break
		case Activated:
			state, _ = ActivateAlarmOnErrors(rctx, v)
			break
		case Alarm:
			state, _ = HandleAlarm(rctx, r)
		}
	}
}

func ActivateWhenNoErrors(rtcx RunContext, v Verify) (WatchdogState, error) {

	return rtcx.Call("verify", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := v.Verify()

		if vErr != nil {
			return Preactivated, vErr
		}

		return Activated, nil
	})
}

func ActivateAlarmOnErrors(rtcx RunContext, v Verify) (WatchdogState, error) {

	return rtcx.Call("verify", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := v.Verify()

		if vErr != nil {
			return Alarm, vErr
		}

		return Activated, nil
	})
}

func HandleAlarm(rtcx RunContext, r Recover) (WatchdogState, error) {

	return rtcx.Call("recover", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := r.Recover()

		if vErr != nil {
			return Preactivated, vErr
		}

		return Activated, nil
	})
}
