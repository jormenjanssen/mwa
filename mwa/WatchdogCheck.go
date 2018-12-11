package main

type WatchdogCheck interface {
	Verify() error
}

type WatchdogReset interface {
	Recover() error
}
