package main

import (
	"fmt"
)

type NetworkHealth struct {
	Address string
}

func (nh *NetworkHealth) VerifyConnection() error {
	fmt.Println(fmt.Printf("Invoking ping to: %v", nh.Address))
	_, err := Ping(nh.Address)
	return err
}

func (nh *NetworkHealth) ShouldRetry(count int, err error) bool {
	return (count < 20)
}
