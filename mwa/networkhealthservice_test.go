package main

import (
	"fmt"
	"testing"
	"time"
)

type attempedStub struct {
	SucceedAfter int
	Count        int
}

func (as *attempedStub) Try() error {

	if as.Count >= as.SucceedAfter {
		return nil
	}

	as.Count++

	return fmt.Errorf("To many failed attempts")
}

func TestNetworkHealthService_TryVerifyMultipleAttempts(t *testing.T) {

	stubFirstRunSucces := attempedStub{Count: 0, SucceedAfter: 0}
	stubSecondRunSucces := attempedStub{Count: 0, SucceedAfter: 2}
	stubEveryAttemptFails := attempedStub{Count: 0, SucceedAfter: 32768}

	type args struct {
		f        func() error
		attempts int
		delay    time.Duration
	}
	tests := []struct {
		name    string
		nh      NetworkHealthService
		args    args
		wantErr bool
	}{
		{name: "FirstRunSucces", nh: NetworkHealthService{Address: "127.0.0.1"}, args: args{attempts: 3, f: stubFirstRunSucces.Try}, wantErr: false},
		{name: "SecondRunSucces", nh: NetworkHealthService{Address: "127.0.0.1"}, args: args{attempts: 3, f: stubSecondRunSucces.Try}, wantErr: false},
		{name: "EveryAttemptFails", nh: NetworkHealthService{Address: "127.0.0.1"}, args: args{attempts: 3, f: stubEveryAttemptFails.Try}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.nh.TryVerifyMultipleAttempts(tt.args.f, tt.args.attempts, tt.args.delay); (err != nil) != tt.wantErr {
				t.Errorf("NetworkHealthService.TryVerifyMultipleAttempts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
