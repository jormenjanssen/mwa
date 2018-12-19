package main

import (
	"net"
	"time"

	netlink "github.com/jormenjanssen/netlink"
)

func GetIpv4TargetForAdapterGatewayWithTimeout(adapter string, duration time.Duration) (net.IP, error) {

	startTime := time.Now()

	for time.Since(startTime) < duration {

		ip, _ := GetIpv4TargetForAdapterGateway(adapter)

		if ip != nil {
			return ip, nil
		}
	}

	return nil, TimeOutError
}

func GetIpv4TargetForAdapterGateway(adapter string) (net.IP, error) {

	routes, err := netlink.NetworkGetRoutes()

	if err == nil {

		for _, route := range routes {
			if route.Iface.Name == adapter && route.IPNet.IP.To4() != nil {
				return route.IPNet.IP, nil
			}
		}
	}

	return nil, err
}
