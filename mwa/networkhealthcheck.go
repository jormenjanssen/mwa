package main

type NetworkHealthCheck interface {
	Perform(addr string) error
}
