package main

import (
	"fmt"
	"testing"
)

type testStub struct {
	Count      int
	MaxFail    int
	MaxRetry   int
	RetryCount int
}

func (t *testStub) Test() error {
	if t.Count > t.MaxFail {
		return fmt.Errorf("Max fail calls reached")
	}

	t.Count++
	return nil
}

func (t *testStub) HandleError(count int, err error) bool {
	if t.MaxRetry < count {
		return true
	}

	return false
}

func TestKeepAlive(t *testing.T) {

	type args struct {
		f func() error
		r func(int, error) bool
	}

	stub := testStub{MaxFail: 10, MaxRetry: 3}
	arg := args{f: stub.Test,
		r: stub.HandleError}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "BasicTest", args: arg, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := KeepAlive(tt.args.f, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("KeepAlive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
