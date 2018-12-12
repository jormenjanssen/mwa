package main

import (
	"fmt"
	"testing"
)

type alwaysSuccesTestStub struct{}

func (s alwaysSuccesTestStub) Verify() error {
	return nil
}

type ActivateFailRecoverTestStub struct {
	state *WatchdogState
}

func (s ActivateFailRecoverTestStub) Verify() error {

	if *s.state == Preactivated {
		*s.state = Activated
		return nil
	}

	if *s.state == Activated {
		return fmt.Errorf("Expeced fail")
	}

	if *s.state == Alarm {
		return fmt.Errorf("Expeced fail")
	}

	return nil
}

func TestWatchdog(t *testing.T) {

	type args struct {
		rctx RunContext
		v    Verify
		r    Recover
	}

	preActivateBecomesActivatedArgs := args{
		rctx: NewUnitTestContext(),
		v:    alwaysSuccesTestStub{}}

	cycleArgs := args{
		rctx: NewUnitTestContext(),
		v:    ActivateFailRecoverTestStub{}}

	tests := []struct {
		name string
		args args
	}{
		{name: "PreActivateBecomesActivated", args: preActivateBecomesActivatedArgs},
		{name: "Cycle", args: cycleArgs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Watchdog(tt.args.rctx, tt.args.v, tt.args.r)
		})
	}
}
