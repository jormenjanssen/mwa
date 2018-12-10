package main

import (
	"fmt"
	"time"
)

type NetworkHealth struct {
	Address    string
	Delay      time.Duration
	RetryDelay time.Duration
	Retries    int
}

func (nh *NetworkHealth) VerifyConnection() error {
	fmt.Println(fmt.Printf("Invoking ping to: %v", nh.Address))
	_, err := Ping(nh.Address)

	if err == nil {
		<-time.After(nh.Delay)
	}

	return err
}

func (nh *NetworkHealth) ShouldRetry(count int, err error) bool {

	if count < nh.Retries {
		<-time.After(nh.RetryDelay)
	}

	return (count < nh.Retries)
}
