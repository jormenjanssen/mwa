package main

import "time"

// Get uptime duration
func GetUptime() (time.Duration, error) {
	return getOsUptime()
}
