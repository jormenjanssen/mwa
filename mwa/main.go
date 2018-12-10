package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {

	delay := flag.Duration("delay", 30*time.Second, "The duration to delay after an succesfull attempt")
	retryDelay := flag.Duration("retryDelay", 5*time.Second, "The seconds to delay after an succesfull attempt")
	retries := flag.Int("retries", 3, "The number of retries before giving up")
	host := flag.String("host", "192.168.10.2", "The default host to check for")

	flag.Parse()

	nh := NetworkHealth{Address: *host, Delay: *delay, Retries: *retries, RetryDelay: *retryDelay}

	fmt.Println(fmt.Printf("Running network health against target: %v with interval: %v and max: %v retries (%v delay between failed attempts )", nh.Address, nh.Delay, nh.Retries, nh.RetryDelay))

	for {
		err := KeepAlive(nh.VerifyConnection, nh.ShouldRetry)
		if err != nil {
			fmt.Println(fmt.Printf("Keep-Alive returned error %v", err))
		}
	}
}
