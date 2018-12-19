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

	log.Debugf("Searching for default gateway for adapter: %v", adapter)
	routes, err := netlink.NetworkGetRoutes()

	if err == nil {

		for _, route := range routes {

			if route.Iface != nil {
				log.Debugf("Got route for interface: %v comparing against: %v matching: %v", route.Iface.Name, adapter, route.Iface.Name == adapter)
			}

			if route.Iface != nil && route.Iface.Name == adapter {

				if route.IPNet != nil {
					log.Debugf("found gateway interface IP: %v", route.IPNet.IP)
					log.Debugf("route is default: %v", route.Default)
				}

				if route.IPNet != nil && route.IPNet.IP != nil && route.IPNet.IP.To4() != nil {
					return route.IPNet.IP, nil
				}
			}
		}
	}

	return nil, err
}
