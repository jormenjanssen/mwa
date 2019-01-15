package main

import (
	"net/http"
	"time"
)

type HttpHealthCheck struct {
	Timeout time.Duration
}

func (h HttpHealthCheck) Perform(addr string) error {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	_, err := client.Get(addr)
	return err
}
