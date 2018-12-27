package main

import (
	"context"
	"reflect"
	"testing"
)

func TestStateContext_Call(t *testing.T) {
	type args struct {
		operation string
		f         func(ctx context.Context, ws WatchdogState) (WatchdogState, error)
	}
	tests := []struct {
		name     string
		stateCtx *StateContext
		args     args
		want     WatchdogState
		wantErr  bool
	}{
		{name: "Simple activation test", want: Activated, wantErr: false,
			stateCtx: &StateContext{ctx: context.TODO()}, args: args{operation: "op1", f: func(ctx context.Context, ws WatchdogState) (WatchdogState, error) { return Activated, nil }}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.stateCtx.Call(tt.args.operation, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateContext.Call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateContext.Call() = %v, want %v", got, tt.want)
			}
		})
	}
}
