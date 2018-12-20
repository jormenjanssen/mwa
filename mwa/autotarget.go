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

	byNameInterface, err := net.InterfaceByName(adapter)

	if err != nil {
		return nil, err
	}

	addresses, err := byNameInterface.Addrs()
	if err != nil {
		return nil, err
	}

	routes, err := netlink.NetworkGetRoutes()

	if err == nil {

		for _, route := range routes {

			if route.Iface != nil {
				log.Debugf("Got route for interface: %v comparing against: %v matching: %v", route.Iface.Name, adapter, route.Iface.Name == adapter)
			}

			if route.Iface != nil && route.Iface.Name == adapter {

				if route.IPNet != nil {
					log.Debugf("found gateway interface IP: %v", route.IPNet)
					log.Debugf("route is default: %v", route.Default)
				}

				for _, addr := range addresses {

					ip, ipnet, err := net.ParseCIDR(addr.String())

					if err != nil {
						log.Debugf("Could not convert address: %v to ip/ipnet error: %v", route.IPNet, err)
						continue
					}

					// Check if we're not lookback
					if ip.Equal(net.IPv4zero) {
						log.Debugf("Skipped looback (127.0.0.1)")
						continue
					}

					// Check if we're not ipv6
					if ip.To4() == nil {
						log.Debugf("Skipped %v because it's ipv6", ip)
						continue
					}

					if route.IPNet != nil && route.IPNet.IP != nil && route.IPNet.IP.To4() != nil {
						log.Debugf("Checking if ipnet %v contains %v", ipnet, route.IPNet.IP)
					}

					if route.IPNet != nil && route.IPNet.IP != nil && route.IPNet.IP.To4() != nil && ipnet.Contains(route.IPNet.IP) {
						return route.IPNet.IP, nil
					}
				}
			}
		}
	}

	return nil, err
}
