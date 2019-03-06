package main

type Verify interface {
	Verify(attempts int) error
}
