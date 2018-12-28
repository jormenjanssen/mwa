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
		{name: "CompareMasksGatewayOk", args: args{a: net.IPv4Mask(255, 255, 255, 255), b: net.IPv4Mask(0xff, 0xff, 0xff, 0xff)}, want: true},
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

func TestGetTargetHost(t *testing.T) {
	type args struct {
		host                string
		autoDetectInterface string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Preset host test", args: args{host: "this.should.be.fairly.specific.0755", autoDetectInterface: "eth1"}, want: "this.should.be.fairly.specific.0755", wantErr: false},
		{name: "No host no gateway", args: args{host: "", autoDetectInterface: ""}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTargetHost(tt.args.host, tt.args.autoDetectInterface)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTargetHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTargetHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
