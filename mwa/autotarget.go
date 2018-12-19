package main

import (
	"net"
	"time"

	netlink "github.com/jormenjanssen/netlink"
	log "github.com/sirupsen/logrus"
)

func GetIpv4TargetForAdapterGatewayWithTimeout(adapter string, duration time.Duration) (net.IP, error) {

	startTime := time.Now()

	for time.Since(startTime) < duration {

		ip, err := GetIpv4TargetForAdapterGateway(adapter)

		if err != nil {
			log.Debugf("Failed to detect route error: %v", err)
		}

		if ip != nil {
			return ip, nil
		}

		<-time.After(1 * time.Second)
	}

	return nil, TimeOutError
}

func GetIpv4TargetForAdapterGateway(adapter string) (net.IP, error) {

	routes, err := netlink.NetworkGetRoutes()

	if err == nil {

		for _, route := range routes {

			if err == nil {

				log.Debugf("Detected route for adapter: %v with IP: %v IsIpv4: %v", route.Iface.Name, route.IPNet.IP, route.IPNet.IP.To4() != nil)
				if route.Iface.Name == adapter && route.IPNet.IP.To4() != nil {
					return route.IPNet.IP, nil
				}
			}
		}
	}

	return nil, err
}
