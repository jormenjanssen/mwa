package main

import (
	"fmt"
)

func main() {

	nh := NetworkHealth{Address: "192.168.10.122"}

	for {
		err := KeepAlive(nh.VerifyConnection, nh.ShouldRetry)
		if err != nil {
			fmt.Println(fmt.Printf("Keep-Alive returned error %v", err))
		}
	}
}
