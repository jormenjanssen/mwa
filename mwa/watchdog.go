package main

import (
	"context"
)

type WatchdogState int

const (
	Exit         WatchdogState = -1
	Preactivated WatchdogState = 0
	Activated    WatchdogState = 1
	Alarm        WatchdogState = 2
)

func Watchdog(rctx RunContext, wc WatchdogCheck, wr WatchdogReset) {

	state := Preactivated

	for {

		switch state {

		case Exit:
			return
		case Preactivated:
			state, _ = ActivateWhenNoErrors(rctx, wc)
			break
		case Activated:
			state, _ = ActivateAlarmOnErrors(rctx, wc)
			break
		case Alarm:
			state, _ = HandleAlarm(rctx, wr)
		}
	}
}

func ActivateWhenNoErrors(rtcx RunContext, wc WatchdogCheck) (WatchdogState, error) {

	return rtcx.Call("verify", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := wc.Verify()

		if vErr != nil {
			return Preactivated, vErr
		}

		return Activated, nil
	})
}

func ActivateAlarmOnErrors(rtcx RunContext, wc WatchdogCheck) (WatchdogState, error) {

	return rtcx.Call("verify", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := wc.Verify()

		if vErr != nil {
			return Alarm, vErr
		}

		return Activated, nil
	})
}

func HandleAlarm(rtcx RunContext, wr WatchdogReset) (WatchdogState, error) {

	return rtcx.Call("recover", func(ctx context.Context, ws WatchdogState) (WatchdogState, error) {

		vErr := wr.Recover()

		if vErr != nil {
			return Preactivated, vErr
		}

		return Activated, nil
	})
}
