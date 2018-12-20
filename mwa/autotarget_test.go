package main

import (
	"net"
	"testing"
)

func TestCompareNetworkMasks(t *testing.T) {
	type args struct {
		a net.IPMask
		b net.IPMask
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "CompareMasksGatewayOk", args: args{a: net.IPv4Mask(0, 0, 0, 0), b: net.IPv4Mask(0, 0, 0, 0)}, want: true},
		{name: "CompareMasksGatewayFail", args: args{a: net.IPv4Mask(255, 255, 255, 0), b: net.IPv4Mask(0, 0, 0, 0)}, want: false},
		{name: "CompareMasksNetworkOk", args: args{a: net.IPv4Mask(255, 255, 0, 0), b: net.IPv4Mask(255, 255, 0, 0)}, want: true},
		{name: "CompareMasksNetworkFail", args: args{a: net.IPv4Mask(255, 255, 255, 0), b: net.IPv4Mask(255, 255, 0, 0)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareNetworkMasks(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CompareNetworkMasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
