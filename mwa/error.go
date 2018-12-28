package main

import (
	"fmt"
	"time"
)

// LastErrorFunc wraps a function where the result is always an error
// If the function thats being called is returning nil, then we return the alternate (aerr) error
func LastErrorFunc(f func() error, aerr error) error {

	if f == nil {
		return fmt.Errorf("No recovery action is defined we cannot recover")
	}

	err := f()

	if err == nil && aerr != nil {
		err = aerr
	} else if err == nil {
		err = fmt.Errorf("Alternate error is empty")
	}

	return err
}

func TimeOutError(topic string, d time.Duration) error {
	return fmt.Errorf("Exceeded timeout: %v while executing: %v", d, topic)
}
