package main

import (
	"strings"
	"time"
)

func GetNeworkHealthCheck(addr string) NetworkHealthCheck {

	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return HttpHealthCheck{Timeout: 5 * time.Second}
	}

	return PingHealthCheck{}
}
