package main

import (
	"reflect"
	"testing"
	"time"
)

func TestGetNeworkHealthCheck(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want NetworkHealthCheck
	}{
		{name: "ShouldBeHttp", args: args{addr: "http://google.nl"}, want: HttpHealthCheck{Timeout: 5 * time.Second}},
		{name: "ShouldBeHttp", args: args{addr: "https://google.nl"}, want: HttpHealthCheck{Timeout: 5 * time.Second}},
		{name: "ShouldBeIcmp", args: args{addr: "127.0.0.1"}, want: PingHealthCheck{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNeworkHealthCheck(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNeworkHealthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
